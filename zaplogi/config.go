package zaplogi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/lohvht/logfeller"
	"github.com/pkg/errors"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig encapsulates the initialisation of the zap logger
type LogConfig struct {
	// ConsoleLog determines if you want to log to the console.
	ConsoleLog bool `json:"console_log" yaml:"console-log"`
	// LogFileConfigs contain the various rotational file configurations
	LogFileConfigs []LogFileConfig `json:"log_file_configs" yaml:"log-file-configs"`
}

// LogFileConfig is the configuration for the log file
// If unmarshalling directly from JSON/YAML, take note
type LogFileConfig struct {
	// LoggerName is the logger's name to log to. If specified, It will only
	// log to this file when the logger's name is equal to LoggerName.
	// This gives fine control over which files that you want to log to.
	LoggerName string
	// LogRange is the level range to log under. If not specified,
	// defaul to [InfoLevel, InfoLevel]
	LogRange [2]Level
	// Type determines what type of log file to use. If specified,
	// accepts "lumberjack" or "logfeller" as the 2 main log file handler configs
	// to marshal as attributes via the `file_handler` field.
	// Check the respective docs on usage via JSON/YAML marshalling
	// 	- `https://pkg.go.dev/github.com/lohvht/logfeller@v1.0.0`
	// 	- `https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2`
	// Otherwise, no io.Writer will be configured. Omit this field if you would
	// like to set the writer manually.
	Type LogFileType
	io.Writer
}

// logFileConfigJSON is the actual struct to marshal JSON to.
type logFileConfigJSON struct {
	LoggerName  string          `json:"logger_name"`
	LogRange    [2]Level        `json:"log_range"`
	Type        LogFileType     `json:"type"`
	FileHandler json.RawMessage `json:"file_handler"`
}

func (c *LogFileConfig) UnmarshalJSON(data []byte) error {
	var lfc logFileConfigJSON
	err := json.Unmarshal(data, &lfc)
	if err != nil {
		return err
	}
	c.LoggerName = lfc.LoggerName
	c.LogRange = lfc.LogRange
	c.Type = lfc.Type
	switch c.Type {
	case NoWriter:
		return nil
	case Lumberjack:
		if c.Writer != nil {
			return fmt.Errorf("writer was already set for log file config; loggername=%q, logrange=%s, type=%q", c.LoggerName, c.LogRange, c.Type)
		}
		c.Writer = &lumberjack.Logger{}
		return json.Unmarshal(lfc.FileHandler, &c.Writer)
	case Logfeller:
		if c.Writer != nil {
			return fmt.Errorf("writer was already set for log file config; loggername=%q, logrange=%s, type=%q", c.LoggerName, c.LogRange, c.Type)
		}
		c.Writer = &logfeller.File{}
		return json.Unmarshal(lfc.FileHandler, &c.Writer)
	default:
		return fmt.Errorf("invalid type: %q", c.Type)
	}
}

// logFileConfigYAMLBase is the actual struct to marshal the base YAML to.
type logFileConfigYAMLBase struct {
	LoggerName string      `yaml:"logger-name"`
	LogRange   [2]Level    `yaml:"log-range"`
	Type       LogFileType `yaml:"type"`
}

type logFileConfigFileHandlerLumberjack struct {
	FileHandler lumberjack.Logger `yaml:"file-handler"`
}

type logFileConfigFileHandlerLogfeller struct {
	FileHandler logfeller.File `yaml:"file-handler"`
}

func (c *LogFileConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var lfc logFileConfigYAMLBase
	err := unmarshal(&lfc)
	if err != nil {
		return err
	}
	c.LoggerName = lfc.LoggerName
	c.LogRange = lfc.LogRange
	c.Type = lfc.Type
	switch c.Type {
	case NoWriter:
		return nil
	case Lumberjack:
		if c.Writer != nil {
			return fmt.Errorf("writer was already set for log file config; loggername=%q, logrange=%s, type=%q", c.LoggerName, c.LogRange, c.Type)
		}
		var lj logFileConfigFileHandlerLumberjack
		err := unmarshal(&lj)
		if err != nil {
			return err
		}
		c.Writer = &lj.FileHandler
	case Logfeller:
		if c.Writer != nil {
			return fmt.Errorf("writer was already set for log file config; loggername=%q, logrange=%s, type=%q", c.LoggerName, c.LogRange, c.Type)
		}
		var lf logFileConfigFileHandlerLogfeller
		err := unmarshal(&lf)
		if err != nil {
			return err
		}
		c.Writer = &lf.FileHandler
	default:
		return fmt.Errorf("invalid type: %q", c.Type)
	}
	return nil
}

type LogFileType int

const (
	NoWriter LogFileType = iota
	Lumberjack
	Logfeller
)

func (t LogFileType) MarshalText() ([]byte, error) { return []byte(t.String()), nil }

// UnmarshalText unmarshals text to a log file type.
// In particular, this makes it easy to configure log file types using YAML,
// TOML, or JSON files.
func (t *LogFileType) UnmarshalText(text []byte) error {
	if t == nil {
		return errors.New("can't unmarshal a nil *LogFileType")
	}
	if !t.unmarshalText(text) && !t.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognised LogFileType: %q", text)
	}
	return nil
}

func (t *LogFileType) unmarshalText(text []byte) bool {
	switch string(text) {
	case "Lumberjack", "lumberjack", "lj":
		*t = Lumberjack
	case "Logfeller", "logfeller", "lf":
		*t = Logfeller
	case "":
		*t = NoWriter
	default:
		return false
	}
	return true
}

// String returns a lower-case ASCII representation of the log file type
func (t LogFileType) String() string {
	switch t {
	case Lumberjack:
		return "lumberjack"
	case Logfeller:
		return "logfeller"
	case NoWriter:
		return "noWriter"
	default:
		return fmt.Sprintf("LogFileType(%d)", t)
	}
}
