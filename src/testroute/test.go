package selector

import (
	"fmt"
	"middlewares"
	"modules"
	"mr"
	"net/http"

	"github.com/labstack/echo"
	"selector"
	//"filter"
	"filter"
)

type selectController struct {
}

func filterNonApp(c *selector.Context, in mr.AdData) bool {
	return in.CpType == 0
}

type size []int

func (tc *selectController) Select(c echo.Context) error {
	rd := c.Get("RequestData").(*middlewares.RequestData)
	wd := c.Get("WebsiteData").(*mr.WebsiteData)
	size:=c.Get("RequestSize").([]int)

	//call context
	m := selector.Context{
		RequestData: *rd,
		WebsiteData: *wd,
		Size: size,

	}
	x := selector.Apply(&m, selector.GetAdData(), selector.Mix(filter.CheckForSize,filter.CheckWhiteList), 3)
	fmt.Println(len(x))
	return c.JSON(http.StatusOK, x)
}

func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Get("/select", tc.Select)
}

func init() {

	modules.Register(&selectController{})
}
