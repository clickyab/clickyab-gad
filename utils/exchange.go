package utils

import "fmt"

var (
	supplierMap = map[string]string{}
	supplier    = map[string]int64{}
)

// GetSupplier return the id supplier
func GetSupplier(name string) (string, int64, error) {
	if newName, ok := supplierMap[name]; ok {
		name = newName
	}

	if s := supplier[name]; s > 0 {
		return name, s, nil
	}

	return "", 0, fmt.Errorf("supplier %s is not valid", name)
}
