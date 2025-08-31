package controller

import (
	"gpm/common/res"

	"github.com/gin-gonic/gin"
)

type UserLoginReq struct {
	Email    string `json:"email" binding:"required,email,min=5,max=100"`
	Password string `json:"password" binding:"required,min=5,max=16"`
}

func GpmHealthView(c *gin.Context) {
	res.SuccessWithData(c, struct{}{})
}
