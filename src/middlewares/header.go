package middlewares

import (
	"fmt"
	"net/url"

	"gopkg.in/labstack/echo.v3"
)

// Header set no cache header for routes
func Header(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Request().Header.Set("Pragma", "no-cache")
		ctx.Request().Header.Set("Cache-Control", "no-cache")
		var site string = ctx.Request().Header.Get("Origin")
		if site == "" {
			ref, err := url.Parse(ctx.Request().Referer())
			if err == nil {
				site = fmt.Sprintf("%s://%s", ref.Scheme, ref.Host)
			}
		}
		ctx.Request().Header.Set("Access-Control-Allow-Origin", site)
		ctx.Request().Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Request().Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		ctx.Request().Header.Set("Access-Control-Allow-Credentials", "true")

		return next(ctx)
	}
}
