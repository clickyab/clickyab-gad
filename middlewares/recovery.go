package middlewares

import (
	"github.com/clickyab/services/safe"
	"gopkg.in/labstack/echo.v3"
)

// Recovery is the middleware to prevent the panic to crash the app
func Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var res error
		safe.Routine(func() { res = next(ctx) }, ctx.Request())
		return res
	}
}
