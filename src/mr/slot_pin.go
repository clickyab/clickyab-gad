package mr

import (
	"time"
)

// SlotPinData structure to define fix slot
type SlotPinData struct {
	AdData

	Direct   bool      `json:"direct" db:"direct"`
	Chance   int       `json:"chance" db:"chance"`
	Bid      int64     `json:"bid" db:"bid"`
	SlotSize int       `json:"-" db:"slot_size"`
	SlotID   int64     `json:"-" db:"slot_id"`
	Start    time.Time `json:"start" db:"start"`
	End      time.Time `json:"end" db:"end"`
}

// SlotPin slot pin table model
type SlotPin struct {
	ID        int       `json:"id" db:"id"`
	SlotID    int       `json:"slot_id" db:"slot_id"`
	Chance    int       `json:"chance" db:"chance"`
	AdID      int       `json:"ad_id" db:"ad_id"`
	Bid       int64     `json:"bid" db:"bid"`
	Direct    bool      `json:"direct" db:"direct"`
	Start     time.Time `json:"start" db:"start"`
	End       time.Time `json:"end" db:"end"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
