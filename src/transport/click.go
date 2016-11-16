package transport

import (
	"net"
	"time"
)

// Impression is the record for the single impression
type Click struct {
	ID           int64     `json:"id"`
	CopID        int64     `json:"copid"`
	IP           net.IP    `json:"ip"`
	AdID         int64     `json:"adid"`
	SlotID       int64     `json:"sid"`
	CampaignID   int64     `json:"cid"`
	UserAgent    string    `json:"ua"`
	WinnerBid    int64     `json:"winner_bid"`
	InTime       time.Time `json:"it"`
	OutTime      time.Time `json:"ot"`
	FraudReason  int       `json:"fre"`
	SlaID        int64     `json:"sla_id"`
	ImpID        int64     `json:"imp_id"`
	OS           int64     `json:"os"`
	Status       int64     `json:"status"`
	CampaignAdID int64     `json:"cpadid"`
	Rand         string    `json:"rand"`
	TrueView     bool      `json:"tv"`
	// TODO : better status
	Web *WebSiteImp `json:"web"`
}

func (c Click) Validate() bool {
	if c.CopID == 0 ||
		len(c.IP.String()) == 0 ||
		c.SlotID <= 0 ||
		c.AdID <= 0 ||
		c.CampaignID <= 0 ||

		c.Web == nil ||
		c.Web.WebsiteID <= 0 {
		return false
	}
	return true
}
