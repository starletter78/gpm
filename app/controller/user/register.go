package user

import (
	"gpm/app/model"
	"gpm/common/res"
	"gpm/common/util"
	"gpm/global"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRegisterReq struct {
	Email    string `json:"email" binding:"required,email,min=5,max=100"`
	Password string `json:"password" binding:"required,min=5,max=16"`
}

func (UserApi) UserRegisterView(c *gin.Context) {
	var cr UserRegisterReq
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailValid(c, err.Error())
		return
	}
	var user model.User
	global.DB.WithContext(c.Request.Context()).Where("email = ?", cr.Email).Find(&user)
	if user.ID != "" {
		res.FailWithMsg(c, "账号已存在")
		return
	}
	salt := uuid.New().String()
	hash := util.Md5([]byte(cr.Password + salt))
	err := global.DB.WithContext(c.Request.Context()).Create(&model.User{
		Email:    cr.Email,
		Password: hash,
		Salt:     salt,
	}).Error
	if err != nil {
		res.FailWithMsg(c, err.Error())
		return
	}
	res.FailWithMsg(c, "注册成功")
}
