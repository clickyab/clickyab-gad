package middlewares

import (
	"fmt"
	"net/url"

	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

// Header set no cache header for routes
func Header(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set("Pragma", "no-cache")
		ctx.Response().Header().Set("Cache-Control", "no-cache")
		var site string = ctx.Request().Header.Get("Origin")
		if site == "" {
			ref, err := url.Parse(ctx.Request().Referer())
			if err == nil {
				site = fmt.Sprintf("%s://%s", ref.Scheme, ref.Host)
			}
		}
		logrus.Debug(site)
		ctx.Response().Header().Set("Access-Control-Allow-Origin", site)
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")

		return next(ctx)
	}
}
