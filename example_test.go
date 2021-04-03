package logi_test

import (
	"encoding/json"
	"fmt"

	"github.com/lohvht/logfeller"
	"github.com/lohvht/logi/zaplogi"
	"gopkg.in/yaml.v2"
)

func ExampleLogConfig_jsonUnmarshal_logfeller() {
	b := []byte(`{
		"console_log": true,
		"log_file_configs": [
			{
				"log_range": ["info", "fatal"],
				"type": "lf",
				"file_handler": {
					"filename": "info.log",
					"when": "d",
					"rotation_schedule": ["0000:00"],
					"use_local": true,
					"backups": 30
				}
			},
			{
				"log_range": ["warn", "max"],
				"type": "logfeller",
				"file_handler": {
					"filename": "error.log",
					"when": "d",
					"rotation_schedule": ["0000:00"],
					"use_local": true,
					"backups": 30
				}
			},
			{
				"logger_name": "db",
				"log_range": ["debug", "fatal"],
				"type": "Logfeller",
				"file_handler": {
					"filename": "db.log",
					"when": "d",
					"rotation_schedule": ["0000:00"],
					"use_local": true,
					"backups": 30
				}
			},
			{
				"logger_name": "data",
				"type": "lf",
				"file_handler": {
					"filename": "data.log",
					"when": "d",
					"rotation_schedule": ["0000:00"],
					"use_local": true,
					"backups": 30
				}
			}
		]
	}`)
	var logConfig zaplogi.LogConfig
	err := json.Unmarshal(b, &logConfig)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Println(logConfig.ConsoleLog)
	for i, rc := range logConfig.LogFileConfigs {
		fmt.Printf("==================%d of %d==================\n", i+1, len(logConfig.LogFileConfigs))
		fmt.Println(rc.LoggerName)
		fmt.Println(rc.LogRange)
		fmt.Println(rc.Type)
		fmt.Println(rc.Writer.(*logfeller.File).Filename)
		fmt.Println(rc.Writer.(*logfeller.File).When)
		fmt.Println(rc.Writer.(*logfeller.File).RotationSchedule)
		fmt.Println(rc.Writer.(*logfeller.File).UseLocal)
		fmt.Println(rc.Writer.(*logfeller.File).Backups)
	}
	// Output:
	// true
	// ==================1 of 4==================
	//
	// [info fatal]
	// logfeller
	// info.log
	// d
	// [0000:00]
	// true
	// 30
	// ==================2 of 4==================
	//
	// [warn fatal]
	// logfeller
	// error.log
	// d
	// [0000:00]
	// true
	// 30
	// ==================3 of 4==================
	// db
	// [debug fatal]
	// logfeller
	// db.log
	// d
	// [0000:00]
	// true
	// 30
	// ==================4 of 4==================
	// data
	// [info info]
	// logfeller
	// data.log
	// d
	// [0000:00]
	// true
	// 30
}

func ExampleLogConfig_yamlUnmarshal_logfeller() {
	b := []byte(`console-log: yes
log-file-configs:
- log-range: ['info', 'fatal']
  type: lf
  file-handler:
    filename: 'info.log'
    when: 'd'
    rotation-schedule: ['0000:00']
    use-local: yes
    backups: 30
- log-range: ['warn', 'max']
  type: logfeller
  file-handler:
    filename: 'error.log'
    when: 'd'
    rotation-schedule: ['0000:00']
    use-local: yes
    backups: 30
- logger-name: 'db'
  log-range: ['debug', 'fatal']
  type: Logfeller
  file-handler:
    filename: 'db.log'
    when: 'd'
    rotation-schedule: ['0000:00']
    use-local: yes
    backups: 30
- logger-name: 'data'
  type: lf
  file-handler:
    filename: 'data.log'
    when: 'd'
    rotation-schedule: ['0000:00']
    use-local: yes
    backups: 30
`)
	var logConfig zaplogi.LogConfig
	err := yaml.Unmarshal(b, &logConfig)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Println(logConfig.ConsoleLog)
	for i, rc := range logConfig.LogFileConfigs {
		fmt.Printf("==================%d of %d==================\n", i+1, len(logConfig.LogFileConfigs))
		fmt.Println(rc.LoggerName)
		fmt.Println(rc.LogRange)
		fmt.Println(rc.Type)
		fmt.Println(rc.Writer.(*logfeller.File).Filename)
		fmt.Println(rc.Writer.(*logfeller.File).When)
		fmt.Println(rc.Writer.(*logfeller.File).RotationSchedule)
		fmt.Println(rc.Writer.(*logfeller.File).UseLocal)
		fmt.Println(rc.Writer.(*logfeller.File).Backups)
	}
	// Output:
	// true
	// ==================1 of 4==================
	//
	// [info fatal]
	// logfeller
	// info.log
	// d
	// [0000:00]
	// true
	// 30
	// ==================2 of 4==================
	//
	// [warn fatal]
	// logfeller
	// error.log
	// d
	// [0000:00]
	// true
	// 30
	// ==================3 of 4==================
	// db
	// [debug fatal]
	// logfeller
	// db.log
	// d
	// [0000:00]
	// true
	// 30
	// ==================4 of 4==================
	// data
	// [info info]
	// logfeller
	// data.log
	// d
	// [0000:00]
	// true
	// 30
}
