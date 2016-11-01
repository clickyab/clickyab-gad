package transport

import (
	"database/sql"
	"net"
	"time"
)

// Impression is the record for the single impression
type Click struct {
	User           string         `json:"u"`
	IP             net.IP         `json:"ip"`
	AdID           int64          `json:"adid"`
	SlotID         int64          `json:"sid"`
	CampaignID     int64          `json:"cid"`
	UserAgent      string         `json:"ua"`
	InTime         time.Time      `json:"it"`
	OutTime        time.Time      `json:"ot"`
	WinnerBID      int64          `json:"wb"`
	RefererAddress sql.NullString `json:"referrer"`
	ParentUrl      sql.NullString `json:"parent"`
	FraudReason    int            `json:"fre"`
	// TODO : better status
	Status int64       `json:"s"`
	Cookie bool        `json:"c"`
	Imp    *Impression `json:"imp"`
	Web    *WebSiteImp `json:"web"`
}

func (c Click) Validate() bool {
	if len(c.User) == 0 ||
		len(c.IP.String()) == 0 ||
		c.SlotID <= 0 ||
		c.AdID <= 0 ||
		c.CampaignID <= 0 ||
		c.WinnerBID <= 0 ||
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
