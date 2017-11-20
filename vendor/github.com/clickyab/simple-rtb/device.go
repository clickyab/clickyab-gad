package srtb

// Device shows the clients device details
type Device struct {
	UA string `json:"ua"`
	IP string `json:"ip"`
	// TODO : Move Geo in higher level, one of the meaning of the simple is to be in one level only
	Geo      Geo    `json:"geo,omitempty"`
	ConnType int    `json:"connectiontype,omitempty"`
	Carrier  string `json:"carrier,omitempty"`
	Lang     string `json:"lang,omitempty"`
	LAC      string `json:"lac,omitempty"`
	MNC      string `json:"mnc,omitempty"`
	MCC      string `json:"mcc,omitempty"`
	CID      string `json:"cid,omitempty"`
}

// Country is the country object
type Country struct {
	Valid bool   `json:"valid"`
	Name  string `json:"name"`
	// Country using ISO 3166-1 Alpha 3
	ISO string `json:"iso"`
}

// Region of the request
type Region struct {
	Valid bool   `json:"valid"`
	Name  string `json:"name"`
	// Region using ISO 3166-2
	ISO string `json:"iso"`
}

// LatLon is the latitude longitude
type LatLon struct {
	Valid bool    `json:"valid"`
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
}

// ISP network holder
type ISP struct {
	Valid bool   `json:"valid"`
	Name  string `json:"name"`
}

// Geo srtb geo type
type Geo struct {
	LatLon  LatLon  `json:"latlon,omitempty"`
	Region  Region  `json:"region,omitempty"`
	Country Country `json:"country,omitempty"`
	ISP     ISP     `json:"isp,omitempty"`
}
