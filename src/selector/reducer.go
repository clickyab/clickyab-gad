package selector

import (
	"mr"
	"sync"
	"fmt"
	"github.com/labstack/echo"
)

// FilterFunc is the type use to filter the
type FilterFunc func(echo.Context, mr.AdData) bool

// ReduceFunc is not needed since the data type is simple and known at compile time
//type ReduceFunc func ([]mr.AdData, []mr.AdData) []mr.AdData

func Apply(ctx echo.Context, in []mr.AdData, ff FilterFunc, cc int) []mr.AdData {
	if cc < 1 {
		cc = 1
	}
	wg := sync.WaitGroup{}
	wg.Add(len(in))
	sem := make(chan struct{}, cc)
	res := make([]mr.AdData, 0, len(in))
	for i := range in {
		go func(j int) {
			sem <- struct{}{}
			defer func() {
				wg.Done()
				<-sem
			}()
			fmt.Println(j)
			if ff(ctx, in[j]) {
				res = append(res, in[j])
			}
		}(i)
	}

	wg.Wait()

	return res
}
