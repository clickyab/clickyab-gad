package modules

import (
	"config"
	"middlewares"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

func Register(c ...Controller) {
	all = append(all, c...)
}

// Initialize the controller
func Initialize(mountPoint string) *echo.Echo {
	once.Do(func() {
		e = echo.New()
		mid := []echo.MiddlewareFunc{middlewares.Recovery, middlewares.Logger, middlewares.RequestCollector}
		if config.Config.CORS {
			mid = append(mid, middleware.CORS())
		}
		e.Use(mid...)
		for i := range all {
			all[i].Routes(e, mountPoint)
		}
		e.SetLogger(NewLogger())
	})
	//engine.SetLogLevel(log.DEBUG)
	return e
}
