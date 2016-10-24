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
	// TODO : better status
	Status     int64       `json:"s"`
	Cookie     bool        `json:"c"`
	Suspicious bool        `json:"sp"`
	Imp        *Impression `json:"imp"`
	Web        *WebSiteImp `json:"web"`
}
