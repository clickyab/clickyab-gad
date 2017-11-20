package middlewares

import (
	"net/http"
	"time"

	echo "gopkg.in/labstack/echo.v3"

	"fmt"

	"github.com/sirupsen/logrus"
)

// Logger is the middleware for log system
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Start timer
		start := time.Now()
		fields := logrus.Fields{
			"method": c.Request().Method,
			"path":   c.Request().URL.Path,
			"query":  c.Request().URL.String(),
		}

		c.Set("LOG_DATA", fields)
		// Process request
		err := next(c)

		fields = c.Get("LOG_DATA").(logrus.Fields)
		latency := time.Since(start)
		statusCode := c.Response().Status

		fields["latency"] = latency.String()
		fields["status"] = statusCode

		logrus.WithFields(fields).Info(http.StatusText(statusCode))

		return err
	}
}

// SetData sets log data containing time
func SetData(c echo.Context, key string, value interface{}) {
	l, ok := c.Get("LOG_DATA").(logrus.Fields)
	if !ok {
		return // WHY?
	}

	var data interface{}
	switch t := value.(type) {
	case int, int32, int64, float64, float32, string, bool:
		data = t
	case error:
		data = t.Error()
	case time.Duration:
		data = t.Seconds()
	case time.Time:
		data = t.Format(time.RFC3339)
	case fmt.Stringer:
		data = t.String()
	default:
		return // can not accept other types
	}

	l[key] = data
	c.Set("LOG_DATA", l)
}
