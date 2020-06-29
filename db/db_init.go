package db

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// 读取mysql的配置, 初始化mysql
type DbConf struct {
	AliasName  string `json:"alias_name"`
	Host       string `json:"host"`
	MasterHost string `json:"master_host"`
	User       string `json:"user"`
	Password   string `json:"password"`
	DbName     string `json:"db_name"`
}

func GetMasterAliasName(AliasName string) string {
	return AliasName + "_m"
}

func (d *DbConf) GetMasterHost() string {
	if d.MasterHost == "" {
		return d.Host
	}
	return d.MasterHost
}

type Log interface {
	Write(p []byte) (n int, err error)
}

// init mysql params(30, 500,int64(10*time.Minute))
func InitDB(aliasName, user, password, host, dbName string, debugLog bool, dueTimeMs int64, log Log, params ...int64) error {
	orm.Debug = debugLog
	orm.ExecuteTime = time.Duration(dueTimeMs) * time.Millisecond
	orm.DebugLog = orm.NewLog(log)
	if debugLog == false {
		orm.OnlyPrintFail = true
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&loc=Local", user, password, host, dbName)
	return orm.RegisterDataBase(aliasName, "mysql", source, params...)
}

func InitDBAndMaster(conf DbConf, debugLog bool, dueTimeMs int64, log Log, params ...int64) error {
	err := InitDB(conf.AliasName, conf.User, conf.Password, conf.Host, conf.DbName, debugLog, dueTimeMs, log, params...)
	if err != nil {
		return err
	}
	//添加主库
	err = InitDB(GetMasterAliasName(conf.AliasName), conf.User, conf.Password, conf.GetMasterHost(), conf.DbName, debugLog, dueTimeMs, log, params...)
	return err
}
