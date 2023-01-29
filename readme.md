# Logi is a common logging interface library

Logi is a batteries included logging interface library. It uses a wrapped
[zap](https://github.com/uber-go/zap) sugared logger that logs only to stderr/stdout
by default however it can be extended and customised as needed as well.

## Usage

Logi is already setup by default to log to stdin and stderr depending
on the level which you are logging at.
```
package main

import "github.com/lohvht/logi"

func main(){
    logi.Get().Info("Hello World!")
}
```

Logi also provides a way to customise the underlying zap logger as well.
For example, if you would like to add in rotational logging to a file `info.log`
every day at 12am by default via JSON, you can do the following:
```
package main

import (
    "encoding/json"

    "github.com/lohvht/logi"
    "github.com/lohvht/logi/zaplogi"
)

func main() {
    // The JSON config to use to log to info.log and ensure it rotates at
    // 1200am daily
    // This uses logfeller for the rotational schedule logic. Other option
    // for "type" would be lumberjack
    jsonConfig := []byte(`{
		"console_log": true,
		"root_caller_skip": 1,
		"log_file_configs": [
			{
				"log_range": ["info", "fatal"],
				"type": "logfeller",
				"file_handler": {
					"filename": "info.log",
					"when": "d",
					"rotation_schedule": ["0000:00"],
					"use_local": true,
					"backups": 30
				}
			},
		]
	}`)
    // Get logConfig data first via JSON.
	var logConfig zaplogi.LogConfig
	err := json.Unmarshal(b, &logConfig)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
    // initialise the new logger withe config passed in.
    newLogger := zaplogi.NewWithConfig(logConfig)
    // Register this new Logger as the default logger to use for logi
    // any subsequent logi.Get() will return this new logger instead.
    logi.SetDefault(newLogger)

    // This should also print to info.log as well.
    logi.Get().Info("Hello world!")
}
```
Zaplogi also supports YAML too.

Zaplogi's file logging may be customised further than just using logfeller or
lumberjack as `zaplogi.LogFileConfig` accepts any io.Writer.

For example, if in addition to the options specified above, we want to also include
an extra `error.log` that does not have any rotational logic:
```
package main

import (
    "os"
    "encoding/json"

    "github.com/lohvht/logfeller"
    "github.com/lohvht/logi"
    "github.com/lohvht/logi/zaplogi"
)

func main() {
    f, err := os.OpenFile("error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	logConfig := zaplogi.LogConfig{
		ConsoleLog: true,
		LogFileConfigs: []zaplogi.LogFileConfig{
			{
				LogRange: [2]zaplogi.Level{zaplogi.InfoLevel, zaplogi.FatalLevel},
				Writer: &logfeller.File{
					Filename:         "info.log",
					When:             "d",
					RotationSchedule: []string{"0000:00"},
					UseLocal:         true,
					Backups:          30,
				},
			},
			{
				LogRange: [2]zaplogi.Level{zaplogi.WarnLevel, zaplogi.MaxLevel},
				Writer:   f,
			},
		},
	}
    // initialise the new logger withe config passed in.
    newLogger := zaplogi.NewWithConfig(logConfig)
    // Register this new Logger as the default logger to use for logi
    // any subsequent logi.Get() will return this new logger instead.
    logi.SetDefault(newLogger)

    // This should also print to info.log as well.
    logi.Get().Info("Hello world!")
    // This should print to error.log
    logi.Get().Error("I have run into an error!")
}
```
