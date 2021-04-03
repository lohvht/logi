package iface

// Logger is the standard interface for loggers in toc-common. As long as your
// logging framework satisfies this interface and its contract, you should
// replace the default implementation of this Logger to your desired format via
// logger.SetDefault()
type Logger interface {
	// Debug logs a message with some additional context in debug level. The
	// variadic key-value pairs are treated as they are in With. See With for
	// more information
	Debug(msg string, keysAndValues ...interface{})
	// Debugf uses fmt.Sprintf to log a templated message in debug level.
	Debugf(template string, args ...interface{})
	// Info logs a message with some additional context in info level. The
	// variadic key-value pairs are treated as they are in With. See With for
	// more information
	Info(msg string, keysAndValues ...interface{})
	// Infof uses fmt.Sprintf to log a templated message  in info level.
	Infof(template string, args ...interface{})
	// Warn logs a message with some additional context in warn level. The
	// variadic key-value pairs are treated as they are in With. See With for
	// more information
	Warn(msg string, keysAndValues ...interface{})
	// Warnf uses fmt.Sprintf to log a templated message in warn level.
	Warnf(template string, args ...interface{})
	// Error logs a message with some additional context in error level. The
	// variadic key-value pairs are treated as they are in With. See With for
	// more information
	Error(msg string, keysAndValues ...interface{})
	// Errorf uses fmt.Sprintf to log a templated message in error level.
	Errorf(template string, args ...interface{})
	// Panic logs a message with some additional context in panic level and then
	// proceeds to panic. The variadic key-value pairs are treated as they are in
	// With. See With for more information
	Panic(msg string, keysAndValues ...interface{})
	// Panicf uses fmt.Sprintf to log a templated message in panic level and then
	// proceeds to panic.
	Panicf(template string, args ...interface{})
	// Fatal logs a message with some additional context in fatal level and then
	// exits. The variadic key-value pairs are treated as they are in With.
	// See With for more information
	Fatal(msg string, keysAndValues ...interface{})
	// Fatalf uses fmt.Sprintf to log a templated message in fatal level and then
	// exits.
	Fatalf(template string, args ...interface{})

	// With returns a logger that provides additional context to the logger.
	// The arguments passed in should be should be in the order of a key-value
	// pair. e.g.
	// With("arg1", arg1, "arg2", arg2)
	// where arg1 and arg2 are the values that we are interested in.
	With(args ...interface{}) Logger

	// Named returns a new logger with the given name
	Named(loggerName string) Logger
}
