package selector

import (
	"clickyab.com/gad/models"

	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/utils"
)

// Context is the context used in reducer functions
type Context struct {
	middlewares.RequestData
	// TODO : its better to have a unique size array
	Size         map[string]int
	Website      *models.Website
	Province     int64
	ISP          int64
	App          *models.App
	PhoneData    *models.PhoneData
	CellLocation *models.CellLocation
	Campaign     int64
	SlotPins     []models.SlotPinData

	MinBidPercentage float64
}

// FilterFunc is the type use to filter the
type FilterFunc func(*Context, models.AdData) bool

// Mix try to mix multiple filter to single function so there is no need to
// call Apply more than once
func Mix(f ...FilterFunc) FilterFunc {
	return func(c *Context, a models.AdData) bool {
		for i := range f {
			if !f[i](c, a) {
				return false
			}
		}
		return true
	}
}

// Apply get the data and then call filter on each of them concurrently, the
// result is the accepted items
func Apply(ctx *Context, in []models.AdData, ff FilterFunc) map[int][]*models.AdData {
	m := make(map[int][]*models.AdData)
	for i := range in {
		if ff(ctx, in[i]) {
			n := in[i]
			n.WinnerBid = 0
			if n.AdType == utils.AdTypeVideo {
				for _, j := range utils.GetVideoSize() {
					m[j] = append(m[j], &n)
				}
			} else {
				m[in[i].AdSize] = append(m[in[i].AdSize], &n)
			}

		}
	}
	return m
}
