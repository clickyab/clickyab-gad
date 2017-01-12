package transport

import "time"

// Warning is the structure to process at the logger
type Warning struct {
	Level   string
	When    time.Time
	Where   string
	Request []byte
	Message string
}

// GetTopic of the Warning struct
func (Warning) GetTopic() string {
	return "cy.warn"
}

// GetQueue of the Warning
func (Warning) GetQueue() string {
	return "cy_warn_queue"
}
