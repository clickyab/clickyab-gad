package selector

import (
	"fmt"
	"middlewares"
	"modules"
	"mr"
	"net/http"

	"github.com/labstack/echo"
)

type selectController struct {
}

func filterNonApp(c *Context, in mr.AdData) bool {
	return in.CpType == 0
}

func filterSize(c *Context, in mr.AdData) bool {
	return true
}

func (tc *selectController) Select(c echo.Context) error {
	rd := c.Get("RequestData").(*middlewares.RequestData)
	wd := c.Get("WebsiteData").(*mr.WebsiteData)
	//call context
	m := Context{
		RequestData: *rd,
		WebsiteData: *wd,
	}
	x := Apply(&m, GetAdData(), Mix(filterNonApp, filterSize), 3)
	fmt.Println(len(x))
	return c.JSON(http.StatusOK, x)
}

func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Get("/select", tc.Select)
}

func init() {

	modules.Register(&selectController{})
}
