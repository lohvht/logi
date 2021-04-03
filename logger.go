package logi

import (
	"github.com/lohvht/logi/iface"
	"github.com/lohvht/logi/zaplogi"
)

var defaultLogger iface.Logger = zaplogi.NewConsole()

// SetDefault replaces the default logger used by logi
func SetDefault(l iface.Logger) { defaultLogger = l }

// Get returns logi's current default logger
// By default, the logger uses zap and logs only to Stdout and Stderr, you will
// need to set SetDefault to change the logger
func Get() iface.Logger { return defaultLogger }
