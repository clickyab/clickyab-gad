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
	dbmap *gorp.DbMap
	db    *sql.DB
	once  = sync.Once{}
	all   []utils.Initializer
)

type gorpLogger struct {
}

func (g gorpLogger) Printf(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

// Initialize the modules, its safe to call this as many time as you want.
func Initialize() {
	once.Do(func() {
		var err error
		db, err = sql.Open("mysql", config.Config.Mysql.DSN)
		assert.Nil(err)

		db.SetMaxIdleConns(config.Config.Mysql.MaxIdleConnection)
		db.SetMaxOpenConns(config.Config.Mysql.MaxConnection)
		err = db.Ping()
		assert.Nil(err)

		dbmap = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

		if config.Config.DevelMode {
			logger := gorpLogger{}
			dbmap.TraceOn("[db]", logger)
		} else {
			dbmap.TraceOff()
		}
		common.Initialize(db, dbmap)
		for i := range all {
			all[i].Initialize()

		}
	})
}

// Register a new initializer module
func Register(m ...utils.Initializer) {
	all = append(all, m...)
}
