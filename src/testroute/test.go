package testroute

import (
	"modules"

	"github.com/labstack/echo"
)

type testController struct {
}

func (tc *testController) Example(c echo.Context) error {
	c.HTML(200, "OK")

	return nil
}

func (tc *testController) Routes(e *echo.Echo, _ string) {
	e.Get("/test", tc.Example)
}

func init() {
	modules.Register(&testController{})
}
