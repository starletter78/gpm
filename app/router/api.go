package router

import (
	"gpm/app/controller"

	"github.com/gin-gonic/gin"
)

func ApiRoute(r *gin.RouterGroup) {
	app := controller.AdminApi{}.ApiApi
	userRoute := r.Group("api")
	//userRoute.GET("", middleware.JwtMiddleware, middleware.CasbinMiddleware, app.ApiListView)
	userRoute.GET("", app.ApiListView)
}
