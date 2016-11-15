package transport

// Conversion is the conversion tracking system
type Conversion struct {
	ConvID   string `json:"cid"`
	ActionID string `json:"action_id"`
	//ImpID    string `json:"imp_id"`
}
