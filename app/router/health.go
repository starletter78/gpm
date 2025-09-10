package router

import (
	"gpm/app/controller"
	"gpm/app/middleware"

	"github.com/gin-gonic/gin"
)

func HealthRoute(r *gin.RouterGroup) {
	app := controller.GpmApi{}.HealthApi
	healthRoute := r.Group("health")
	healthRoute.Use(middleware.LogMiddleware)
	healthRoute.GET("", app.GpmHealthView)
}
