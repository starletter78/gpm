package api

import (
	"gpm/app/model"
	"gpm/common"
	"gpm/common/res"

	"github.com/gin-gonic/gin"
)

func (ApiApi) ApiOptionsView(c *gin.Context) {
	var cr common.PageInfo
	if err := c.ShouldBindQuery(&cr); err != nil {
		res.FailWithError(c, err)
		return
	}
	result, count, err := common.NewQueryBuilder(
		&model.Api{},
		common.Options{
			PageInfo: cr,
		},
	).Build().GetResult()
	if err != nil {
		res.FailWithError(c, err)
		return
	}
	var _result []model.OptionsRes

	for _, v := range result {
		_result = append(_result, model.OptionsRes{
			Id:   v.ID,
			Name: v.Name,
		})
	}

	res.SuccessWithList(c, _result, count)
}
