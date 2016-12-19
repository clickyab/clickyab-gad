package middlewares

import (
	"config"
	"fmt"
	"net/http"
	"runtime/debug"
	"utils"

	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

// Recovery is the middleware to prevent the panic to crash the app
func Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				ctx.JSON(
					http.StatusInternalServerError,
					struct {
						Error string `json:"error"`
					}{
						Error: http.StatusText(http.StatusInternalServerError),
					},
				)
				stack := debug.Stack()
				dump := "TODO: create dump" //httputil.DumpRequest(ctx.Request(), true)
				data := fmt.Sprintf("Request : \n %s \n\nStack : \n %s", dump, stack)
				logrus.WithField("error", err).Warn(err, data)
				if config.Config.Redmine.Active {
					go utils.RedmineDoError(err, []byte(data))
				}

				if config.Config.Slack.Active {
					go utils.SlackDoMessage(err, ":shit:", utils.SlackAttachment{Text: data, Color: "#AA3939"})
				}
			}
		}()

		return next(ctx)
	}
}
