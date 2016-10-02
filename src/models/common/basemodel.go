package common

import (
	"database/sql"
	"errors"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql" // Make sure postgres is included in any build
	"gopkg.in/gorp.v1"
)

var (
	dbmap *gorp.DbMap
	db    *sql.DB
)

// Manager is a base manager for transaction model
type Manager struct {
	tx    *gorp.Transaction
	dbmap *gorp.DbMap
	db    *sql.DB

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
	m.tx, err = m.dbmap.Begin()
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
	if m.dbmap == nil {
		m.dbmap = dbmap
	}
}

// GetDbMap is for getting the current dbmap
func (m *Manager) GetDbMap() gorp.SqlExecutor {
	if m.transaction {
		return m.tx
	}
	m.sureDbMap()
	return m.dbmap
}

// GetSQLDB return the raw connection to database
func (m *Manager) GetSQLDB() *sql.DB {
	if m.db == nil {
		m.db = db
	}

	return m.db
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
	return m.dbmap.AddTable(i)
}

// AddTableWithName has the same behavior as AddTable, but sets
// table.TableName to name.
func (m *Manager) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	m.sureDbMap()
	return m.dbmap.AddTableWithName(i, name)
}

// AddTableWithNameAndSchema has the same behavior as AddTable, but sets
// table.TableName to name.
func (m *Manager) AddTableWithNameAndSchema(i interface{}, schema string, name string) *gorp.TableMap {
	m.sureDbMap()
	return m.dbmap.AddTableWithNameAndSchema(i, schema, name)
}

// TruncateTables try to truncate tables , useful for tests
func (m *Manager) TruncateTables(tbl string) error {
	m.sureDbMap()
	q := "TRUNCATE " + tbl
	_, err := m.dbmap.Exec(q)
	return err
}

// Initialize the module
func Initialize(d *sql.DB, dbm *gorp.DbMap) {
	dbmap = dbm
	db = d
}
