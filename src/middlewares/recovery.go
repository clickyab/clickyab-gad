package middlewares

import (
	"config"
	"fmt"
	"net/http"
	"runtime/debug"
	"utils"

	"net/http/httputil"

	"github.com/sirupsen/logrus"
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
				dump, _ := httputil.DumpRequest(ctx.Request(), true)
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

// SafeGO run a function in safe manner
func SafeGO(ctx echo.Context, exit bool, continuous bool, f func()) {
	go func() {
		s := make(chan struct{})
		for {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						stack := debug.Stack()
						var dump []byte
						if ctx != nil {
							dump, _ = httputil.DumpRequest(ctx.Request(), true)
						}
						data := fmt.Sprintf("Do a restart on %v \n Request : \n %s \n\nStack : \n %s", exit, dump, stack)
						logrus.WithField("error", err).Warn(err, data)
						if config.Config.Redmine.Active {
							go utils.RedmineDoError(err, []byte(data))
						}

						if config.Config.Slack.Active {
							go utils.SlackDoMessage(err, ":shit:", utils.SlackAttachment{Text: data, Color: "#AA3939"})
						}

						if exit {
							logrus.Fatal(err)
						}
					}
					s <- struct{}{} // allow to run once more
				}()
				f()
			}()
			<-s // block it here until the defered function is done
			if !continuous {
				break
			}
		}
	}()
}
