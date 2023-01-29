// zaplogi provides a zaplogger interface replacement
package zaplogi

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lohvht/logfeller"
	"github.com/lohvht/logi/iface"
)

type Logger struct {
	zaplog *zap.SugaredLogger
}

// defaultEncoderConfig returns the default encoding used. Note that EncodeLevel
// is encoded in colour by default.
func defaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey: "msg",
		TimeKey:    "timestamp",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000Z0700"))
		},
		LevelKey:       "level",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		NameKey:        "logger",
		CallerKey:      "caller",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		StacktraceKey:  "stacktrace",
		EncodeDuration: zapcore.SecondsDurationEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
	}
}

// NewWithConfig returns a Logger with the given config. Logger
// implemen logger.Logger interface, so you may use SetDefault to replace the
// default logger.
// With this function, you can customise the log output and how you would want
// to initialise the logger. However, usually New will suffice.
func NewWithConfig(c LogConfig) (*Logger, error) {
	encConf := defaultEncoderConfig()
	var enc zapcore.Encoder
	var childCores []zapcore.Core
	options := []zap.Option{zap.AddCallerSkip(c.RootCallerSkip), zap.AddCaller()}
	if c.ConsoleLog {
		enc = zapcore.NewConsoleEncoder(encConf)
		stdoutPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			// log debugs to stdout
			return lvl < zapcore.WarnLevel
		})
		stdErrPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.WarnLevel
		})
		stdoutCore := zapcore.NewCore(enc, zapcore.Lock(os.Stdout), stdoutPriority)
		stderrCore := zapcore.NewCore(enc, zapcore.Lock(os.Stderr), stdErrPriority)
		childCores = append(childCores, stdoutCore, stderrCore)
		// add in zap.Development()
		options = append(options, zap.Development())
	}
	var loggerNamesToExclude []string
	// Collect all the logger names to exclude first
	for _, logConf := range c.LogFileConfigs {
		if logConf.LoggerName != "" {
			loggerNamesToExclude = append(loggerNamesToExclude, logConf.LoggerName)
		}
	}
	var Errs []error
	// change the encoding back
	encConf.EncodeLevel = zapcore.CapitalLevelEncoder
	enc = zapcore.NewConsoleEncoder(encConf)
	for _, logConf := range c.LogFileConfigs {
		low, high := logConf.LogRange[0], logConf.LogRange[1]
		if low > high {
			Errs = append(Errs, fmt.Errorf("log level high (%s) is smaller than low (%s)", high.String(), low.String()))
			continue
		}
		lvlFn := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			// Dont log debugs at all
			return lvl >= zapcore.Level(low) && lvl <= zapcore.Level(high)
		})
		if logConf.Writer != nil {
			// Only allow logging if the writer is initialised.
			var childCore zapcore.Core
			if logConf.LoggerName != "" {
				childCore = newExclusiveCore([]string{logConf.LoggerName}, true, zapcore.NewCore(enc, zapcore.AddSync(logConf), lvlFn))
			} else {
				childCore = newExclusiveCore(loggerNamesToExclude, false, zapcore.NewCore(enc, zapcore.AddSync(logConf), lvlFn))
			}
			childCores = append(childCores, childCore)
		}
	}
	if len(Errs) > 0 {
		var buf bytes.Buffer
		buf.WriteString("logger: errors initialising logs,")
		for _, err := range Errs {
			buf.WriteString(fmt.Sprintf("[%s]", err.Error()))
		}
		return nil, errors.New(buf.String())
	}
	core := zapcore.NewTee(childCores...)
	logger := zap.New(core, options...).Sugar()
	zl := &Logger{zaplog: logger}
	defer func() {
		innerErr := logger.Sync()
		if innerErr != nil {
			var pathErr *os.PathError
			if errors.As(innerErr, &pathErr) && strings.HasPrefix(pathErr.Path, "/dev/std") {
				// Ignore logging any sync related errors for logging to
				// /dev/std* as this may pollute logging to console.
				return
			}
			// only log other errors to debug. the logger should still be usable
			zl.Debug("logger syncing error", "syncerror", innerErr.Error())
		}
	}()
	return zl, nil
}

// NewConsole returns a logger that will only log to console (stdout and
// stderr). NewConsole implement logger.Logger and this logger is the
// initial logger for the package.
func NewConsole() *Logger {
	l, _ := NewWithConfig(LogConfig{ConsoleLog: true, RootCallerSkip: 1})
	return l
}

// NewDefault creates 2 log files with max backups that rotates every day
// at 12am.
func NewDefault(logDir string, backups int, logConsole bool) *Logger {
	l, _ := NewWithConfig(LogConfig{
		ConsoleLog:     logConsole,
		RootCallerSkip: 1,
		LogFileConfigs: []LogFileConfig{
			{
				LogRange: [2]Level{InfoLevel, MaxLevel},
				Writer: &logfeller.File{
					Filename:         filepath.Join(logDir, "info.log"),
					When:             "d",
					RotationSchedule: []string{"0000:00"},
					Backups:          backups,
					UseLocal:         true,
				},
			},
			{
				LogRange: [2]Level{WarnLevel, MaxLevel},
				Writer: &logfeller.File{
					Filename:         filepath.Join(logDir, "error.log"),
					When:             "d",
					RotationSchedule: []string{"0000:00"},
					Backups:          backups,
					UseLocal:         true,
				},
			},
		},
	})
	return l
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	if l == nil {
		return
	}
	l.zaplog.Debugw(msg, keysAndValues...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zaplog.Debugf(template, args...)
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	if l == nil {
		return
	}
	l.zaplog.Infow(msg, keysAndValues...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zaplog.Infof(template, args...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	if l == nil {
		return
	}
	l.zaplog.Warnw(msg, keysAndValues...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zaplog.Warnf(template, args...)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	if l == nil {
		return
	}
	l.zaplog.Errorw(msg, keysAndValues...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zaplog.Errorf(template, args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	if l == nil {
		panic(errors.New(fmt.Sprintf(template, args...)))
	}
	l.zaplog.Panicf(template, args...)
}

func (l *Logger) Panic(msg string, keysAndValues ...interface{}) {
	if l == nil {
		kvs := append([]interface{}{msg}, keysAndValues...)
		panic(errors.New(fmt.Sprintln(kvs...)))
	}
	l.zaplog.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	if l == nil {
		kvs := append([]interface{}{msg}, keysAndValues...)
		fmt.Fprintln(os.Stderr, kvs...)
		os.Exit(1)
	}
	l.zaplog.Fatalw(msg, keysAndValues...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	if l == nil {
		fmt.Fprintf(os.Stderr, template, args...)
		os.Exit(1)
	}
	l.zaplog.Fatalf(template, args...)
}

func (l *Logger) With(args ...interface{}) iface.Logger {
	newLogger := l.zaplog.With(args...)
	return &Logger{zaplog: newLogger}
}

func (l *Logger) Named(loggerName string) iface.Logger {
	newLogger := l.zaplog.Named(loggerName)
	return &Logger{zaplog: newLogger}
}

func (l *Logger) CallSkip(skips int) iface.Logger {
	newLogger := l.zaplog.WithOptions(zap.AddCallerSkip(skips))
	return &Logger{zaplog: newLogger}
}

// exclusiveCore is a wrapper around zapcore.Core. It takes a list of logger
// names and check entries depending on the include parameter:
//   - If include==true: check entries if entry's logger name is in the
//     list of logger names passed in
//   - else: check entries if entry's logger name is not in the list of logger
//     name passed in
type exclusiveCore struct {
	loggerNames []string
	include     bool
	zapcore.Core
}

func newExclusiveCore(loggerNames []string, include bool, core zapcore.Core) zapcore.Core {
	// loggerNames are sorted only once upon creation of a new core.
	sort.Strings(loggerNames)
	return &exclusiveCore{
		loggerNames: loggerNames,
		include:     include,
		Core:        core,
	}
}

func (c *exclusiveCore) With(fields []zapcore.Field) zapcore.Core {
	return &exclusiveCore{
		loggerNames: c.loggerNames,
		include:     c.include,
		Core:        c.Core.With(fields),
	}
}

// Check overrides the underlying zapcore.Core implementation by having a check on whether to include
// based on the entry's logger name.
// nolint // to satisfy zapcore.Core interface
func (c *exclusiveCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	i := sort.SearchStrings(c.loggerNames, ent.LoggerName)
	checkCond := i < len(c.loggerNames) && c.loggerNames[i] == ent.LoggerName
	if !c.include {
		// if not include ==> exclude if found in logger names.
		checkCond = !checkCond
	}
	if checkCond {
		return c.Core.Check(ent, ce)
	}
	return ce
}
