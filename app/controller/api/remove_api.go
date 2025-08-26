package api

import (
	"fmt"
	"gpm/app/model"
	"gpm/common/res"
	"gpm/global"

	"github.com/gin-gonic/gin"
)

func (ApiApi) RemoveApiView(c *gin.Context) {
	var cr model.IdListReq
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(c, err)
		return
	}
	if err := global.DB.WithContext(c.Request.Context()).Delete(&model.Api{}, cr.IdList).Error; err != nil {
		res.FailWithError(c, err)
		return
	}
	res.SuccessWithMsg(c, fmt.Sprintf("删除成功%d条", len(cr.IdList)))
}
