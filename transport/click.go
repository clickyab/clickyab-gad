package transport

import (
	"net"
	"time"
)

// Click is the record for the single Click
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
	SLAID        int64     `json:"sla_id"`
	ImpID        int64     `json:"imp_id"`
	OS           int64     `json:"os"`
	Status       int64     `json:"status"`
	CampaignAdID int64     `json:"cpadid"`
	Rand         string    `json:"rand"`
	TrueView     bool      `json:"tv"`
	// TODO : better status
	Web    *WebSiteImp `json:"web"`
	App    *AppImp     `json:"app"`
	AdSize int         `json:"ad_size"`
}

// Validate try to validate the click in
func (c Click) Validate() bool {
	if c.CopID == 0 ||
		len(c.IP.String()) == 0 ||
		c.SlotID <= 0 ||
		c.AdID <= 0 ||
		c.CampaignID <= 0 ||

		(c.Web == nil || c.Web.WebsiteID <= 0) && (c.App == nil || c.App.AppID <= 0) {
		return false
	}
	return true
}

// GetTopic return the click topic
func (c Click) GetTopic() string {
	return "cy.click"
}

// GetQueue return the click queue name to use
func (c Click) GetQueue() string {
	return "cy_click_queue"
}
