package flags

import (
	"context"
	"gpm/app/model"
	"gpm/global"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
)

func FlagsDb() {
	sqlDB, err := global.DB.DB()
	if err != nil {
		logrus.WithContext(context.Background()).Fatal(err.Error())
		return
	}
	defer func() {
		err = sqlDB.Close()
		if err != nil {
			logrus.WithContext(context.Background()).Fatal(err.Error())
			return
		}
	}()
	err = global.DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Tenant{},
		&model.ActionLog{},
		&model.Tenant{},
		&model.Menu{},
		&gormadapter.CasbinRule{},
		&model.Api{},
		&model.Doc{},
		&model.DocDir{},
		&model.UserBlack{},
		&model.TokenBlack{},
	)
	if err != nil {
		logrus.Fatal(err)
		return
	}
}
