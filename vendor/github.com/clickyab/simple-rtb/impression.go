package srtb

// Impression each bid request contains some impression which has data about requested ad
type Impression struct {
	ID       string  `json:"id"`
	Banner   *Banner `json:"banner"`
	BidFloor float64 `json:"bidfloor"`
	Secure   int     `json:"secure"`
}
