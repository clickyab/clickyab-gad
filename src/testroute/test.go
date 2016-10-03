package testroute

import (
	"modules"

	"github.com/labstack/echo"
	"mr"
)

type testController struct {
}

func (tc *testController) Example(c echo.Context) error {
	m := mr.NewManager()

	x, err := m.LoadAds()
	if err != nil {
		return c.HTML(500, err.Error())
	}

	return c.JSON(200, x)
}

func (tc *testController) Routes(e *echo.Echo, _ string) {
	e.Get("/test", tc.Example)
}

func init() {

	modules.Register(&testController{})
}
