package db

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Log interface {
	Write(p []byte) (n int, err error)
}

// init mysql params(30, 500,int64(10*time.Minute))
func InitDB(aliasName, user, password, host, dbName string, debugLog bool, log Log, params ...int64) error {
	orm.Debug = debugLog
	orm.DebugLog = orm.NewLog(log)
	if debugLog == false {
		orm.OnlyPrintFail = true
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Local", user, password, host, dbName)
	return orm.RegisterDataBase(aliasName, "mysql", source, params...)
}
