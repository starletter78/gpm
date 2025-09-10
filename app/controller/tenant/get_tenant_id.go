package tenant

import (
	"gpm/app/model"
	"gpm/common/res"
	"gpm/global"

	"github.com/gin-gonic/gin"
)

type GetTenantIdReq struct {
	TenantName string `json:"tenantName"`
}

func (TenantApi) GetTenantIdView(c *gin.Context) {
	var cr GetTenantIdReq
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(c, err)
		return
	}
	var tenant model.Tenant
	err := global.DB.WithContext(c.Request.Context()).Find(&tenant, "name = ?", cr.TenantName).Error
	if err != nil {
		res.FailWithError(c, err)
		return
	}
	res.SuccessWithData(c, gin.H{
		"tenantId": tenant.ID,
	})
}
