package health

import (
	"gpm/common/res"

	"github.com/gin-gonic/gin"
)

func (HealthApi) GpmHealthView(c *gin.Context) {
	res.SuccessWithData(c, struct{}{})
}
