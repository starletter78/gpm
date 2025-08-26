package api

import (
	"gpm/app/model"
	"gpm/common"
	"gpm/common/res"

	"github.com/gin-gonic/gin"
)

type ApiListReq struct {
	common.PageInfo
}

func (ApiApi) ApiListView(c *gin.Context) {
	var cr ApiListReq
	if err := c.ShouldBindQuery(&cr); err != nil {
		res.FailWithError(c, err)
		return
	}
	result, count, err := common.NewQueryBuilder(model.Api{}, common.Options{
		PageInfo: cr.PageInfo,
		Context:  c.Request.Context(),
	}).Build().GetResult()
	if err != nil {
		res.FailWithError(c, err)
		return
	}
	res.SuccessWithList(c, result, count)
}
