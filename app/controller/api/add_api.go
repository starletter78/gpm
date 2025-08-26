package api

import (
	"gpm/app/model"
	"gpm/common/res"
	"gpm/global"

	"github.com/gin-gonic/gin"
)

type AddApiReq struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Method   string `json:"method"`
	TenantID string `json:"tenantId"`
	Auth     bool   `json:"auth"`
	Status   bool   `json:"status"`
	MenuID   string `json:"menuId"`
}

func (ApiApi) AddApiView(c *gin.Context) {
	var cr AddApiReq
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(c, err)
		return
	}
	var api = model.Api{
		Name:     cr.Name,
		Path:     cr.Path,
		Method:   cr.Method,
		TenantID: cr.TenantID,
		Auth:     cr.Auth,
		Status:   cr.Status,
		MenuID:   cr.MenuID,
	}
	err := global.DB.WithContext(c.Request.Context()).Create(&api).Error
	if err != nil {
		res.FailWithError(c, err)
		return
	}
	res.SuccessWithMsg(c, "添加成功")
}
