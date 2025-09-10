package router

import (
	"gpm/app/controller"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.RouterGroup) {
	app := controller.GpmApi{}.UserApi
	userRoute := r.Group("user")
	userRoute.GET("login", app.UserLoginView)
	userRoute.POST("register", app.UserRegisterView)
}
