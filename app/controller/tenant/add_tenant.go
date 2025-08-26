package tenant

import (
	"gpm/app/model"
	"gpm/common/res"
	"gpm/global"

	"github.com/gin-gonic/gin"
)

type AddApiReq struct {
	Name string `json:"name"`
}

func (TenantApi) AddTenantView(c *gin.Context) {
	var cr AddApiReq
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(c, err)
		return
	}
	var tenant = model.Tenant{
		Name: cr.Name,
	}
	err := global.DB.WithContext(c.Request.Context()).Create(&tenant).Error
	if err != nil {
		res.FailWithError(c, err)
		return
	}
	res.SuccessWithMsg(c, "添加成功")
}
