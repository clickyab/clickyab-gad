package transport

const (
	// Delimiter is the delimiter in all redis key
	Delimiter = "_"
	// User is for user keys
	User = "U"
	// Campaign for campaign keys
	Campaign = "CP"
	// Advertise for advertise keys
	Advertise = "AD"
	// Slot for slot keys
	Slot = "S"
	// Website is for website keys
	Website = "W"
	// CustomClickURL is for app keys
	CustomClickURL = "CCU"
	// CustomClickParam is for app keys
	CustomClickParam = "CCP"
	// CustomClickType is for app keys
	CustomClickType = "CCT"

	// FraudPrefix for all fraud key
	FraudPrefix = "F"
	// ImpSubKey is for impression key inside each master hash
	ImpSubKey = "I"
	// ClickSubKey is for click count
	ClickSubKey = "C"

	// CappingKey is the key for capping system
	CappingKey = "CAP2"
	// ImpKey is the key for impression
	ImpKey = "ImpKey"
	// ClickKey is the key for click
	ClickKey = "CLK"
	// MegaKey is the key for mega
	MegaKey = "MGA"
)
