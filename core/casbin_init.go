package core

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
	"gpm/global"
)

func InitCasbin() *casbin.Enforcer {
	a, err := gormadapter.NewAdapterByDBUseTableName(global.DB, "", "casbin_rule")
	if err != nil {
		logrus.Fatal(err)
	}
	e, err := casbin.NewEnforcer("conf/rbac_with_domains_model.conf", a)
	if err != nil {
		logrus.Fatal(err)
	}

	// Or you can use an existing DB "abc" like this:
	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	// a := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/abc", true)

	// Load the policy from DB.
	err = e.LoadPolicy()
	if err != nil {
		logrus.Fatal(err)
	}
	return e
}
