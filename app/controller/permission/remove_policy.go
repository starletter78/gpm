package permission

import (
	"gpm/common/res"
	"gpm/global"

	"github.com/gin-gonic/gin"
)

// 给用户或者角色添加权限
type RemovePolicyReq struct {
	SubId   string `json:"subId" binding:"required"`
	SubType string `json:"subType" binding:"required,oneof=user role" `
	ObjId   string `json:"objId" binding:"required"`
	ObjType string `json:"objType" binding:"required,oneof=api doc menu"`
	Action  string `json:"action" binding:"required,oneof=get post put delete read write owen"`
}

// 添加权限给用户或者角色
func (PermissionApi) RemovePolicyView(c *gin.Context) {
	var cr AddPolicyReq
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(c, err)
		return
	}
	tenant := c.GetString("tenant")
	sub := cr.SubType + ":" + cr.SubId
	obj := cr.ObjType + ":" + cr.ObjId
	_, err = global.CasbinEnforcer.RemovePolicy(sub, tenant, obj, cr.Action)
	if err != nil {
		res.FailWithError(c, err)
		return
	}
	res.SuccessWithMsg(c, "权限添加成功")
}
