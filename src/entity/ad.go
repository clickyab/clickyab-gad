package entity

// AdType is the
type AdType string

const (
	// AdTypeBanner is the banner type
	AdTypeBanner AdType = "banner"
	// AdTypeDynamic is the dynamic type. the code is html
	AdTypeDynamic AdType = "dyn"
	// AdTypeVideo is the video type
	AdTypeVideo AdType = "video"
)

// Advertise is the single advertise interface
type Advertise interface {
	// GetID return the id of advertise
	ID() int64

	Type() AdType

	Campaign() Campaign

	Capping() Capping

	SetCapping(Capping)

	SetCPM(int64)

	CPM() int64

	SetWinnerBID(int64)

	WinnerBID() int64
	// AdCTR the ad ctr from database (its not calculated from
	AdCTR() float64
	// SetCTR set the calculated CTR
	SetCTR(float64)
	// CTR get the calculated CTR
	CTR() float64
}
