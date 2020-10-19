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

func curLogger() *logs.BeeLogger {
	if DefaultLog == nil {
		InitLog("dev")
	}
	return DefaultLog
}

func Emergency(format string, v ...interface{}) {
	curLogger().Emergency(format, v...)
}

func Alert(format string, v ...interface{}) {
	curLogger().Alert(format, v...)
}

func Critical(format string, v ...interface{}) {
	curLogger().Critical(format, v...)
}

func Error(format string, v ...interface{}) {
	curLogger().Error(format, v...)
}

func Warning(format string, v ...interface{}) {
	curLogger().Warning(format, v...)
}

func Notice(format string, v ...interface{}) {
	curLogger().Notice(format, v...)
}

func Info(format string, v ...interface{}) {
	curLogger().Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	curLogger().Debug(format, v...)
}
