package modules

import (
	"config"
	"middlewares"
	"sync"

	"gopkg.in/labstack/echo.v3"
)

var (
	all  []Controller
	once = &sync.Once{}
	e    *echo.Echo
)

// Controller is the main interface for controllers
type Controller interface {
	Routes(*echo.Echo, string)
}

// Register function
func Register(c ...Controller) {
	all = append(all, c...)
}

// Initialize the controller
func Initialize(mountPoint string) *echo.Echo {
	once.Do(func() {
		e = echo.New()
		mid := []echo.MiddlewareFunc{middlewares.Recovery, middlewares.CORS(), middlewares.Logger}
		e.Use(mid...)
		for i := range all {
			all[i].Routes(e, mountPoint)
		}
		e.Logger = NewLogger()
	})
	//engine.SetLogLevel(log.DEBUG)
	return e
}
