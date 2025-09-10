package api_service

import (
	"gpm/app/controller/api"
	"gpm/app/data"
	"gpm/app/model"

	"github.com/gin-gonic/gin"
)

func (ApiService) AddApiService(c *gin.Context, cr api.AddApiReq) (ApiData *model.Api, err error) {
	return data.NewData().ApiData.AddApiData(c, cr)
}
