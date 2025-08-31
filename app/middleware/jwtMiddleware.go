package middleware

import (
	"context"
	jwt2 "gpm/app/service/jwt"
	"gpm/common/res"

	"github.com/gin-gonic/gin"
)

func JwtMiddleware(c *gin.Context) {
	Authorization := c.GetHeader("Authorization")
	auth := c.GetBool("auth")
	if Authorization == "" && auth {
		res.FailToken(c)
		c.Abort()
		return
	}
	if Authorization == "" {
		res.FailWithMsg(c, "token不存在")
		c.Abort()
		return
	}
	jwt := jwt2.NewJWT()

	accessClaims, err := jwt.ParseAccessToken(Authorization)
	if err != nil {
		res.FailWithError(c, err)
		c.Abort()
		return
	}
	ctx := context.WithValue(c.Request.Context(), "userId", accessClaims.Id)
	c.Request = c.Request.WithContext(ctx)
	c.Set("userId", accessClaims.Id)
	c.Set("user", accessClaims)
}
