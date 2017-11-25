package srtb

// Bid part of every bidResponse
type Bid struct {
	ID       string   `json:"id"`
	AdID     string   `json:"ad_id,omitempty"`
	ImpID    string   `json:"imp_id"`
	Price    float64  `json:"price"`
	AdMarkup string   `json:"adm"`
	Width    int      `json:"w"`
	Height   int      `json:"h"`
	Cat      []string `json:"cat,omitempty"`
	WinURL   string   `json:"win_url,omitempty"`
	BillURL  string   `json:"bill_url,omitempty"`
}
