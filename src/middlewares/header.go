package middlewares

import "gopkg.in/labstack/echo.v3"

// Header set no cache header for routes
func Header(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Request().Header.Set("Pragma", "no-cache")
		ctx.Request().Header.Set("Cache-Control", "no-cache")
		return next(ctx)
	}
}
