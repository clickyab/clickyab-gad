package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

// RemoveSlotPins remove fix slot from ad pool
func RemoveSlotPins(c *selector.Context, in mr.AdData) bool {
	for i := range c.SlotPins {
		if c.SlotPins[i].AdID == in.AdID {
			return false
		}
	}
	return true
}