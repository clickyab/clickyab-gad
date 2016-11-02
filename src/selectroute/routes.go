package selectroute

import "github.com/labstack/echo"


// Routes function @todo
func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Get("/select", tc.Select)
	e.Get("/show/:mega/:ad", tc.Show)
}
