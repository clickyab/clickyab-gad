package selector

import (
	"fmt"
	"modules"
	"mr"
	"net/http"
	"middlewares"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type selectController struct {
}

func filterNonApp(c *Context, in mr.AdData) bool {
	logrus.Info("Hi, its me")
	return in.CpType == 0
}

func filterSize(c *Context, in mr.AdData) bool {
	return true
}

func (tc *selectController) Select(c echo.Context) error {
	RequestData := c.Get("RequestData").(*middlewares.RequestData)
	//call context
	m := Context{
		RequestData: *RequestData,
	}
	x := Apply(m, GetAdData(), Mix(filterNonApp, filterSize), 3)
	fmt.Println(len(x))
	return c.JSON(http.StatusOK, x)
}

func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Get("/select", tc.Select)
}

func init() {

	modules.Register(&selectController{})
}
