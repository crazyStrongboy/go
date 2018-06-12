package db

import (
	"errors"
	"fmt"

	. "github.com/polaris1119/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var MasterDB *xorm.Engine

var dns string

func init() {
	mysqlConfig, err := ConfigFile.GetSection("mysql")
	if err != nil {
		fmt.Println("get mysql config error:", err)
		return
	}

	fillDns(mysqlConfig)

	// 启动时就打开数据库连接
	if err = initEngine(); err != nil {
		panic(err)
	}
}

var (
	ConnectDBErr = errors.New("connect db error")
	UseDBErr     = errors.New("use db error")
)

func fillDns(mysqlConfig map[string]string) {
	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"],
		mysqlConfig["charset"])
}

func initEngine() error {
	var err error

	MasterDB, err = xorm.NewEngine("mysql", dns)
	if err != nil {
		return err
	}

	maxIdle := ConfigFile.MustInt("mysql", "max_idle", 2)
	maxConn := ConfigFile.MustInt("mysql", "max_conn", 10)

	MasterDB.SetMaxIdleConns(maxIdle)
	MasterDB.SetMaxOpenConns(maxConn)

	showSQL := ConfigFile.MustBool("xorm", "show_sql", false)
	logLevel := ConfigFile.MustInt("xorm", "log_level", 1)

	MasterDB.ShowSQL(showSQL)
	MasterDB.Logger().SetLevel(core.LogLevel(logLevel))

	return nil
}
