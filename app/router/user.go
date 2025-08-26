package router

import (
	"github.com/gin-gonic/gin"
	"gpm/app/controller"
	"gpm/app/middleware"
)

func UserRoute(r *gin.RouterGroup) {
	app := controller.AdminApi{}.UserApi
	userRoute := r.Group("user")
	userRoute.GET("login", middleware.JwtMiddleware, middleware.CasbinMiddleware, app.UserLoginView)
	userRoute.POST("register", middleware.JwtMiddleware, middleware.CasbinMiddleware, app.UserRegisterView)
}
