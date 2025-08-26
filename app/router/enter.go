package router

import (
	"gpm/app/middleware"
	"gpm/global"

	"github.com/gin-gonic/gin"
)

func Run() {
	engine := gin.Default()
	r := engine.Group("gpm")
	r.Use(middleware.LogMiddleware, middleware.ArgsCheckMiddleware)
	UserRoute(r)
	SearchRoute(r)
	ApiRoute(r)
	err := engine.Run(global.Config.System.Addr())
	if err != nil {
		return
	}
}
