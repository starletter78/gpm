package middleware

import (
	"gpm/common/res"
	"gpm/global"

	"github.com/gin-gonic/gin"
)

func CasbinMiddleware(c *gin.Context) {
	if c.GetBool("auth") {
		return
	}
	tenant := c.GetHeader("tenant")
	userId := c.GetString("userId")
	role := c.GetStringSlice("role")
	if tenant == "" || userId == "" || role == nil {
		res.FailWithMsg(c, "权限鉴定信息缺失")
		c.Abort()
		return
	}

	ok, err := global.CasbinEnforcer.Enforce()
	if err != nil {
		res.FailWithError(c, err)
		c.Abort()
		return
	}
	if !ok {
		res.FailAuth(c)
		c.Abort()
		return
	}
}
