package log

import (
	"github.com/astaxie/beego/logs"
	"os"
)

var DefaultLog *logs.BeeLogger

func isPathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func Emergency(format string, v ...interface{}) {
	DefaultLog.Emergency(format, v...)
}

func Alert(format string, v ...interface{}) {
	DefaultLog.Alert(format, v...)
}

func Critical(format string, v ...interface{}) {
	DefaultLog.Critical(format, v...)
}

func Error(format string, v ...interface{}) {
	DefaultLog.Error(format, v...)
}

func Warning(format string, v ...interface{}) {
	DefaultLog.Warning(format, v...)
}

func Notice(format string, v ...interface{}) {
	DefaultLog.Notice(format, v...)
}

func Info(format string, v ...interface{}) {
	DefaultLog.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	DefaultLog.Debug(format, v...)
}
