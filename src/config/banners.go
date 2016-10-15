package config

var sizes = map[string]int{
	"120x600":  1,
	"160x600":  2,
	"300x250":  3,
	"336x280":  4,
	"468x60":   5,
	"728x90":   6,
	"120x240":  7,
	"320x50":   8,
	"800x440":  9,
	"300x600":  11,
	"970x90":   12,
	"970x250":  13,
	"250x250":  14,
	"300x1050": 15,
	"320x480":  16,
	"48x320":   17,
	"128x128":  18,
}

// GetSize return the size of a banner in clickyab std
func GetSize(size string) (int, error) {
	for key, value := range sizes {
		if key == size {
			return value, nil
		}
	}

	return 0, nil
}
