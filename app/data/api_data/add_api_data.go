package api_data

import (
	"gpm/app/controller/api"
	"gpm/app/model"
	"gpm/global"

	"github.com/gin-gonic/gin"
)

func (ApiData) AddApiData(c *gin.Context, cr api.AddApiReq) (data *model.Api, err error) {
	var apis = model.Api{
		Name:     cr.Name,
		Path:     cr.Path,
		Method:   cr.Method,
		TenantID: cr.TenantID,
		Auth:     cr.Auth,
		Status:   cr.Status,
		MenuID:   cr.MenuID,
	}
	err = global.DB.WithContext(c.Request.Context()).Create(&apis).Error
	if err != nil {
		return nil, err
	}
	return &apis, nil
}
