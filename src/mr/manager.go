package mr

import (
	"models"
	"models/common"
)

// Manager is the model manager for aaa package
type Manager struct {
	common.Manager
}

// NewManager create and return a manager for this module
func NewManager() *Manager {
	return &Manager{}
}

// Initialize aaa package
func (m *Manager) Initialize() {
	m.AddTableWithName(Slot{}, "slots").SetKeys(true, "ID")
	m.AddTableWithName(CookieProfile{}, "cookie_profiles").SetKeys(true, "ID")
	m.AddTableWithName(App{}, "apps").SetKeys(true, "ID")
	m.AddTableWithName(Website{}, "websites").SetKeys(true, "WID")
}

func init() {
	models.Register(NewManager())
}
