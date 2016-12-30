package common

import (
	"database/sql"
	"errors"

	"config"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql" // Make sure postgres is included in any build
	"gopkg.in/gorp.v2"
)

var (
	rdbmap *gorp.DbMap
	rdb    *sql.DB

	wdbmap *gorp.DbMap
	wdb    *sql.DB
)

// Manager is a base manager for transaction model
type Manager struct {
	tx     *gorp.Transaction
	rdbmap *gorp.DbMap
	rdb    *sql.DB
	wdbmap *gorp.DbMap
	wdb    *sql.DB

	transaction bool
}

// InTransaction return true if this manager s in transaction
func (m *Manager) InTransaction() bool {
	return m.transaction
}

// Begin is for begin transaction
func (m *Manager) Begin() error {
	var err error
	if m.transaction {
		logrus.Panic("already in transaction")
	}
	m.sureDbMap()
	m.tx, err = m.wdbmap.Begin()
	if err == nil {
		m.transaction = true
	}
	return err
}

// Commit is for committing transaction. panic if transaction is not started
func (m *Manager) Commit() error {
	if !m.transaction {
		logrus.Panic("not in transaction")
	}
	err := m.tx.Commit()
	if err != nil {
		return err
	}
	m.tx = nil
	m.transaction = false
	return nil
}

// Rollback is for RollBack transaction. panic if transaction is not started
func (m *Manager) Rollback() error {
	if !m.transaction {
		logrus.Panic("Not in transaction")
	}
	err := m.tx.Rollback()

	if err != nil {
		return err
	}

	m.transaction = false
	return nil
}

func (m *Manager) sureDbMap() {
	if m.rdbmap == nil || m.wdbmap == nil {
		m.rdbmap = rdbmap
		m.wdbmap = wdbmap
	}
}

// GetRDbMap is for getting the current dbmap
func (m *Manager) GetRDbMap() gorp.SqlExecutor {
	if m.transaction {
		return m.tx
	}
	m.sureDbMap()
	return m.rdbmap
}

// GetRSQLDB return the raw connection to database
func (m *Manager) GetRSQLDB() *sql.DB {
	if m.rdb == nil {
		m.rdb = rdb
	}

	return m.rdb
}

// GetWDbMap is for getting the current dbmap
func (m *Manager) GetWDbMap() gorp.SqlExecutor {
	if m.transaction {
		return m.tx
	}
	m.sureDbMap()
	return m.wdbmap
}

// GetWSQLDB return the raw connection to database
func (m *Manager) GetWSQLDB() *sql.DB {
	if m.wdb == nil {
		m.wdb = wdb
	}

	return m.wdb
}

// GetProperDBMap try to get the current writer for development mode
func (m *Manager) GetProperDBMap() gorp.SqlExecutor {
	if config.Config.DevelMode {
		return m.GetWDbMap()
	}
	return m.GetRDbMap()
}

// Hijack try to hijack into a transaction
func (m *Manager) Hijack(ts gorp.SqlExecutor) error {
	if m.transaction {
		return errors.New("already in transaction")
	}
	t, ok := ts.(*gorp.Transaction)
	if !ok {
		return errors.New("there is no transaction to hijack")
	}

	m.transaction = true
	m.tx = t

	return nil
}

// AddTable registers the given interface type with gorp. The table name
// will be given the name of the TypeOf(i).  You must call this function,
// or AddTableWithName, for any struct type you wish to persist with
// the given DbMap.
//
// This operation is idempotent. If i's type is already mapped, the
// existing *TableMap is returned
func (m *Manager) AddTable(i interface{}) *gorp.TableMap {
	m.sureDbMap()
	return m.wdbmap.AddTable(i)
}

// AddTableWithName has the same behavior as AddTable, but sets
// table.TableName to name.
func (m *Manager) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	m.sureDbMap()
	return m.wdbmap.AddTableWithName(i, name)
}

// AddTableWithNameAndSchema has the same behavior as AddTable, but sets
// table.TableName to name.
func (m *Manager) AddTableWithNameAndSchema(i interface{}, schema string, name string) *gorp.TableMap {
	m.sureDbMap()
	return m.wdbmap.AddTableWithNameAndSchema(i, schema, name)
}

// TruncateTables try to truncate tables , useful for tests
func (m *Manager) TruncateTables(tbl string) error {
	m.sureDbMap()
	q := "TRUNCATE " + tbl
	_, err := m.wdbmap.Exec(q)
	return err
}

// Initialize the module
func Initialize(rd *sql.DB, rdbm *gorp.DbMap, wd *sql.DB, wdbm *gorp.DbMap) {
	rdbmap = rdbm
	rdb = rd

	wdbmap = wdbm
	wdb = wd

}
