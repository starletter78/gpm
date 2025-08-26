package router

import (
	"gpm/app/controller"
	"gpm/app/middleware"

	"github.com/gin-gonic/gin"
)

func SearchRoute(r *gin.RouterGroup) {
	app := controller.AdminApi{}.SearchApi.File
	userRoute := r.Group("search")
	userRoute.GET("fileTree", middleware.JwtMiddleware, middleware.CasbinMiddleware, app.FileTreeView)
	userRoute.GET("fileSearch", middleware.JwtMiddleware, middleware.CasbinMiddleware, app.FileSearchView)
}
