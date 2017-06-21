package middlewares

import (
	"config"
	"os"

	"src/gopkg.in/labstack/echo.v3"
)

// ServerID set server id header for routes
func ServerID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sid := config.Config.ServerID
		if sid == "" {
			n, _ := os.Hostname()
			sid = n
		}

		c.Response().Header().Set("X-Server-Id", sid)
		return next(c)
	}
}
