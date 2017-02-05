package models

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"services/mysql"
	"strings"
)

// Manager is the type to handle connection
type Manager struct {
	mysql.Manager
}

// SharpArray type is the hack to handle # splited text in our database
type SharpArray string

// Initialize is the base point to register tables if required
func (m *Manager) Initialize() {

}

// Scan convert the json array ino string slice
func (pa *SharpArray) Scan(src interface{}) error {
	s := &sql.NullString{}
	err := s.Scan(src)
	if err != nil {
		return err
	}

	if s.Valid {
		*pa = SharpArray(s.String)
	} else {
		*pa = ""
	}
	return nil

}

// Value try to get the string slice representation in database
func (pa SharpArray) Value() (driver.Value, error) {
	s := sql.NullString{}
	s.Valid = pa != ""
	s.String = string(pa)

	return s.Value()
}

// Has check exist value in sharpArray
func (pa SharpArray) Has(empty bool, in ...int64) (x bool) {
	if len(in) == 0 || len(pa) == 0 {
		return empty
	}
	if len(in) == 1 {
		return strings.Contains(string(pa), fmt.Sprintf("#%d#", in[0]))
	}

	reg := []string{}
	for i := range in {
		reg = append(reg, fmt.Sprintf("#%d#", in[i]))
	}

	return regexp.MustCompile("(" + strings.Join(reg, "|") + ")").MatchString(string(pa))
}

// Match check for at least one match
func (pa SharpArray) Match(empty bool, in SharpArray) bool {
	if len(in) == 0 || len(pa) == 0 {
		return empty
	}
	inTrim := strings.Trim(string(in), "# \n\t")
	return regexp.MustCompile("(" + strings.Replace(inTrim, "#", "#|#", -1) + ")").MatchString(string(pa))
}

// NewManager create and return a manager for this module
func NewManager() *Manager {
	return &Manager{}
}

func init() {
	mysql.Register(&Manager{})
}
