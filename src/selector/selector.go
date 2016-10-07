package selector

import (
	"fmt"
	"modules"
	"mr"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type selectController struct {
}

func filterNonApp(c echo.Context, in mr.AdData) bool {
	logrus.Info("Hi, its me")
	return in.CpType == 0
}

func filterSize(c echo.Context, in mr.AdData) bool {
	size := c.Get("ccc").(int)
	return in.AdSize == size
}

func (tc *selectController) Select(c echo.Context) error {
	c.Set("ccc", 3)
	x := Apply(c, GetAdData(), Mix(filterNonApp, filterSize), 3)
	fmt.Println(len(x))
	return c.JSON(http.StatusOK, x)
}

func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Get("/select", tc.Select)
}

func init() {

	modules.Register(&selectController{})
}
