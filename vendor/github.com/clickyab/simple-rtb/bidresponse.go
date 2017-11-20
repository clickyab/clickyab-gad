package srtb

// BidResponse is the object that exchange return to publisher
type BidResponse struct {
	ID   string `json:"id"`
	Bids []Bid  `json:"bids"`
}
