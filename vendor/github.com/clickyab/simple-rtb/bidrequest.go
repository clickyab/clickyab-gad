package srtb

// BidRequest is what is accepted from each publisher
type BidRequest struct {
	ID     string       `json:"id"`
	Imp    []Impression `json:"imp"`
	Site   *Site        `json:"site,omitempty"`
	App    *App         `json:"app,omitempty"`
	Device *Device      `json:"device"`
	User   *User        `json:"user,omitempty"`
	Test   int          `json:"test"`
	// TODO : Who thinks this is simple???
	TMax int      `json:"tmax,omitempty"`
	BCat []string `json:"bcat,omitempty"`
}
