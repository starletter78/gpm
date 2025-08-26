package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	jwt2 "gpm/app/service/jwt"
	"gpm/common/res"
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
		return
	}
	jwt := jwt2.NewJWT()

	accessClaims, err := jwt.ParseAccessToken(Authorization)
	if err != nil {
		return
	}
	ctx := context.WithValue(c.Request.Context(), "userId", accessClaims.Id)
	c.Request = c.Request.WithContext(ctx)
	c.Set("user", accessClaims)
}
