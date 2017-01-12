package transport

// Conversion is the conversion tracking system
type Conversion struct {
	ConvID   string `json:"cid"`
	ActionID string `json:"action_id"`
	//ImpID    string `json:"imp_id"`
}

// GetTopic of the Conversion struct
func (Conversion) GetTopic() string {
	return "cy.conv"
}

// GetQueue of the Conversion
func (Conversion) GetQueue() string {
	return "cy_conv_queue"
}
