package global

import (
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
	"gpm/conf"
)

var (
	Config         *conf.Config
	DB             *gorm.DB
	CasbinEnforcer *casbin.Enforcer
)
