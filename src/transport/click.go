package transport

import (
	"database/sql"
	"net"
	"time"
)

// Impression is the record for the single impression
type Click struct {
	CopID          string         `json:"copid"`
	IP             net.IP         `json:"ip"`
	AdID           int64          `json:"adid"`
	SlotID         int64          `json:"sid"`
	CampaignID     int64          `json:"cid"`
	UserAgent      string         `json:"ua"`
	Time           time.Time      `json:"it"`
	RefererAddress sql.NullString `json:"referrer"`
	ParentUrl      sql.NullString `json:"parent"`
	FraudReason    int            `json:"fre"`
	// TODO : better status
	Imp *Impression `json:"imp"`
	Web *WebSiteImp `json:"web"`
}

func (c Click) Validate() bool {
	if len(c.CopID) == 0 ||
		len(c.IP.String()) == 0 ||
		c.SlotID <= 0 ||
		c.AdID <= 0 ||
		c.CampaignID <= 0 ||
		c.Imp.WinnerBID <= 0 ||
		len(c.RefererAddress.String) == 0 ||
		len(c.ParentUrl.String) == 0 ||
		c.Imp == nil ||
		c.Imp.Web == nil ||
		c.Web == nil ||
		c.Web.WebsiteID <= 0 {
		return false
	}
	return true
}
