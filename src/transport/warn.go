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
