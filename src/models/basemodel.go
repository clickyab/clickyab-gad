package models

import (
	"assert"
	"config"
	"database/sql"
	"models/common"
	"sync"
	"utils"

	"github.com/Sirupsen/logrus"
	"gopkg.in/gorp.v1"
)

var (
	rdbmap *gorp.DbMap
	wdbmap *gorp.DbMap
	rdb    *sql.DB
	wdb    *sql.DB
	once   = sync.Once{}
	all    []utils.Initializer
)

type gorpLogger struct {
}

func (g gorpLogger) Printf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

// Initialize the modules, its safe to call this as many time as you want.
func Initialize() {
	once.Do(func() {
		var err error
		rdb, err = sql.Open("mysql", config.Config.Mysql.RDSN)
		assert.Nil(err)

		wdb, err = sql.Open("mysql", config.Config.Mysql.WDSN)
		assert.Nil(err)

		rdb.SetMaxIdleConns(config.Config.Mysql.MaxIdleConnection)
		rdb.SetMaxOpenConns(config.Config.Mysql.MaxConnection)
		wdb.SetMaxIdleConns(config.Config.Mysql.MaxIdleConnection)
		wdb.SetMaxOpenConns(config.Config.Mysql.MaxConnection)

		err = rdb.Ping()
		assert.Nil(err)

		err = wdb.Ping()
		assert.Nil(err)

		rdbmap = &gorp.DbMap{Db: rdb, Dialect: gorp.MySQLDialect{}}
		wdbmap = &gorp.DbMap{Db: wdb, Dialect: gorp.MySQLDialect{}}

		if config.Config.DevelMode {
			logger := gorpLogger{}
			rdbmap.TraceOn("[rdb]", logger)
			wdbmap.TraceOn("[wdb]", logger)
		} else {
			rdbmap.TraceOff()
			wdbmap.TraceOff()
		}
		common.Initialize(rdb, rdbmap, wdb, wdbmap)
		for i := range all {
			all[i].Initialize()

		}
	})
}

// Register a new initializer module
func Register(m ...utils.Initializer) {
	all = append(all, m...)
}
