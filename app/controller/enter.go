package controller

import (
	"gpm/app/controller/api"
	"gpm/app/controller/doc"
	"gpm/app/controller/health"
	"gpm/app/controller/menu"
	"gpm/app/controller/permission"
	"gpm/app/controller/role"
	"gpm/app/controller/search"
	"gpm/app/controller/tenant"
	"gpm/app/controller/user"
)

type GpmApi struct {
	UserApi       user.UserApi
	SearchApi     search.SearchApi
	ApiApi        api.ApiApi
	DocApi        doc.DocApi
	TenantApi     tenant.TenantApi
	RoleApi       role.RoleApi
	MenuApi       menu.MenuApi
	PermissionApi permission.PermissionApi
	HealthApi     health.HealthApi
}
