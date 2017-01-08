package selector

import (
	"mr"

	"config"
	"middlewares"
)

// Context is the context used in reducer functions
type Context struct {
	middlewares.RequestData
	// TODO : its better to have a unique size array
	Size         map[string]int
	Website      *mr.Website
	Province     *mr.Province
	App          *mr.App
	PhoneData    *mr.PhoneData
	CellLocation *mr.CellLocation
	Campaign     int64
}

// FilterFunc is the type use to filter the
type FilterFunc func(*Context, mr.AdData) bool

// Mix try to mix multiple filter to single function so there is no need to
// call Apply more than once
func Mix(f ...FilterFunc) FilterFunc {
	return func(c *Context, a mr.AdData) bool {
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
func Apply(ctx *Context, in []mr.AdData, ff FilterFunc) map[int][]*mr.AdData {
	m := make(map[int][]*mr.AdData)
	for i := range in {
		if ff(ctx, in[i]) {
			n := in[i]
			if n.AdType == config.AdTypeVideo {
				for _, j := range config.GetVideoSize() {
					m[j] = append(m[j], &n)
				}
			} else {
				m[in[i].AdSize] = append(m[in[i].AdSize], &n)
			}

		}
	}
	return m
}
