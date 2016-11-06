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
	m.AddTableWithName(Slots{},"slots").SetKeys(true,"ID")
}

func init() {
	models.Register(NewManager())
}
