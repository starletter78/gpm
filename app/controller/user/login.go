package user

import (
	"github.com/gin-gonic/gin"
	"gpm/app/model"
	"gpm/app/service/jwt"
	"gpm/common/res"
	"gpm/common/util"
	"gpm/global"
)

type UserLoginReq struct {
	Email    string `json:"email" binding:"required,email,min=5,max=100"`
	Password string `json:"password" binding:"required,min=5,max=16"`
}

func (UserApi) UserLoginView(c *gin.Context) {
	var cr UserLoginReq
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailValid(c, err.Error())
		return
	}
	var user model.User
	global.DB.WithContext(c.Request.Context()).Where("email = ?", cr.Email).Find(&user)
	if user.ID == "" {
		res.FailWithMsg(c, "账号不存在")
		return
	}
	hash := util.Md5([]byte(cr.Password + user.Salt))
	if hash != user.Password {
		res.FailWithMsg(c, "密码错误")
		return
	}

	newJwt := jwt.NewJWT()

	pairToken, err := newJwt.GenPairToken(c.GetString("tenant"), user.ID)
	if err != nil {
		res.FailWithMsg(c, err.Error())
		return
	}
	res.SuccessWithData(c, pairToken)
}
