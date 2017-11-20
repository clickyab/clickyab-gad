package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
)

var (
	exchangeSuppliers = config.RegisterString("exchange.supplier", "randomsupkey:supname:1234", "comma separated")
)

// GetSupplier return the id supplier
func GetSupplier(key string) (string, int64, error) {
	a := exchangeSuppliers.String()
	sArr := strings.Split(a, ",")
	if len(sArr) > 0 {
		for i := range sArr {
			eArr := strings.Split(sArr[i], ":")
			if len(eArr) == 3 {
				if eArr[0] == key {
					userID, err := strconv.ParseInt(eArr[2], 10, 64)
					assert.Nil(err)
					return eArr[1], userID, nil
				}
			}
		}
	}
	return "", 0, errors.New("supplier not found")
}
