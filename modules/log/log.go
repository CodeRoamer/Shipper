package log

import (
	"os"

	"github.com/gogits/logs"
)

var (
	loggers   []*logs.BeeLogger
)

func init() {
	NewLogger(0, "console", `{"level": 0}`)
}


func NewLogger(bufLen int64, mode, config string) {
	logger := logs.NewLogger(bufLen)

	isExist := false
	for _, l := range loggers {
		if l.Adapter == mode {
			isExist = true
			l = logger
		}
	}
	if !isExist {
		loggers = append(loggers, logger)
	}
	logger.SetLogFuncCallDepth(3)
	if err := logger.SetLogger(mode, config); err != nil {
		Fatal("Fail to set logger(%s): %v", mode, err)
	}
}

func Trace(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Trace(format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Debug(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Info(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Error(format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Warn(format, v...)
	}
}

func Critical(format string, v ...interface{}) {
	for _, logger := range loggers {
		logger.Critical(format, v...)
	}
}

func Fatal(format string, v ...interface{}) {
	Error(format, v...)
	for _, l := range loggers {
		l.Close()
	}
	os.Exit(2)
}
