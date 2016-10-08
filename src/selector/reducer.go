package selector

import (
	"mr"
	"sync"

	"github.com/labstack/echo"
	"middlewares"
)

type Context struct {
	middlewares.RequestData
	Size []int
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

// ReduceFunc is not needed since the data type is simple and known at compile time
//type ReduceFunc func ([]mr.AdData, []mr.AdData) []mr.AdData

// Apply get the data and then call filter on each of them concurrently, the
// result is the accepted items
func Apply(ctx *Context, in []mr.AdData, ff FilterFunc, cc int) []mr.AdData {
	if cc < 1 {
		cc = 1
	}
	wg := sync.WaitGroup{}
	wg.Add(len(in))
	sem := make(chan struct{}, cc)
	res := make([]mr.AdData, 0, len(in))
	for i := range in {
		sem <- struct{}{}
		go func(j int) {
			defer func() {
				wg.Done()
				<-sem
			}()
			if ff(ctx, in[j]) {
				res = append(res, in[j])
			}
		}(i)
	}

	wg.Wait()

	return res
}
