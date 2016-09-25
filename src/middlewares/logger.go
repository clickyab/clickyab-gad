package middlewares

import (
	"net/http"
	"time"

	"net"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

// Logger is the middleware for log system
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Start timer
		start := time.Now()
		path := c.Request().URL().Path

		// Process request
		err := next(c)

		latency := time.Since(start)

		clientIP := c.Request().RemoteAddress()
		if ip := c.Request().Header().Get(echo.HeaderXRealIP); ip != "" {
			clientIP = ip
		} else if ip = c.Request().Header().Get(echo.HeaderXForwardedFor); ip != "" {
			clientIP = ip
		} else {
			clientIP, _, _ = net.SplitHostPort(clientIP)
		}
		method := c.Request().Method
		statusCode := c.Response().Status()
		logrus.WithFields(
			logrus.Fields{
				"Method":   method,
				"Path":     path,
				"Latency":  latency,
				"ClientIP": clientIP,
				"Status":   statusCode,
				"Err":      err,
			},
		).Info(http.StatusText(statusCode))

		return err
	}
}
