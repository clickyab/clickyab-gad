package config

import "fmt"

var (
	supplier = map[string]int64{
		"adro": 100000,
		"adad": 100001,
		"saba": 100002,
	}
)

// GetSupplier return the id supplier
func GetSupplier(name string) (int64, error) {
	s, ok := supplier[name]
	if ok {
		return s, nil
	}

	return 0, fmt.Errorf("supplier %s is not valid", name)
}
