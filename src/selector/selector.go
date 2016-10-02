package selector

import (
	"modules"

	"github.com/labstack/echo"
)

type selectController struct {
}

func (tc *selectController) Select(c echo.Context) error {
	c.HTML(200, "OK")

	return nil
}

func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Get("/select", tc.Select)
}

func init() {

	modules.Register(&selectController{})
}
