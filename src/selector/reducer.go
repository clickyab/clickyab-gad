package selector

import (
	"mr"

	"middlewares"
)

// Context type @todo
type Context struct {
	middlewares.RequestData
	Size []int
	mr.WebsiteData
	mr.Country2Info
	SlotPublic []string
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
func Apply(ctx *Context, in []mr.AdData, ff FilterFunc, cc int) map[int][]*mr.MinAdData {
	//if cc < 1 {
	//	cc = 1
	//}
	//wg := sync.WaitGroup{}
	//wg.Add(len(in))
	//sem := make(chan struct{}, cc)
	//res := make([]mr.AdData, 0, len(in))
	m := make(map[int][]*mr.MinAdData)
	for i := range in {
		/*sem <- struct{}{}
		go func(j int) {
			defer func() {
				wg.Done()
				<-sem
			}()*/
		if ff(ctx, in[i]) {
			//res = append(res, in[i])
			n := in[i].MinAdData
			m[in[i].AdSize] = append(m[in[i].AdSize], &n)

		}
		/*}(i)
		}

		wg.Wait()*/
	}

	//fmt.Println(res)
	return m
}
