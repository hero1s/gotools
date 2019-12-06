package log

import (
	"github.com/astaxie/beego/logs"
	"os"
	"path/filepath"
)

func InitLog(runMode string) *logs.BeeLogger {
	DefaultLog = logs.NewLogger()
	DefaultLog.EnableFuncCallDepth(true)
	DefaultLog.SetLogFuncCallDepth(3)
	if runMode == "dev" {
		DefaultLog.SetLogger(logs.AdapterConsole)
	} else {
		workPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		logDir := filepath.Join(workPath, "log")
		logFile := filepath.Join(logDir, "log.txt")
		if !isPathExist(logDir) {
			err = os.Mkdir(logDir, 0755)
			if err != nil {
				panic(err)
			}
			f, err := os.Create(logFile)
			if err != nil {
				panic(err)
			}
			f.Close()
		} else {
			f, err := os.OpenFile(logFile,os.O_RDWR|os.O_CREATE|os.O_APPEND,0755)
			if err != nil {
				panic(err)
			}
			f.Close()
		}
		logFile = "log/log.txt"
		configs := `{"filename":"` + logFile + `","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`
		err = DefaultLog.SetLogger(logs.AdapterMultiFile, configs)
		if err != nil {
			panic(err)
		}
	}
	return DefaultLog
}
