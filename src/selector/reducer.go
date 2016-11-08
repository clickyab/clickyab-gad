package selector

import (
	"mr"

	"config"
	"middlewares"
)

// Context type @todo
type Context struct {
	middlewares.RequestData
	// TODO : its better to have a unique size array
	Size    map[string]int
	Website *mr.WebsiteData
	Country *mr.CountryInfo
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
func Apply(ctx *Context, in []mr.AdData, ff FilterFunc) map[int][]*mr.MinAdData {
	m := make(map[int][]*mr.MinAdData)
	for i := range in {
		if ff(ctx, in[i]) {
			n := in[i].MinAdData
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
