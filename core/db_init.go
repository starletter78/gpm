package core

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"gpm/app/service/log"
	"gpm/global"
	"time"
)

func InitDB() *gorm.DB {
	// 创建我们自定义的 GORM Logger
	customGormLogger := log.NewGormLogger()
	if len(global.Config.DB) == 0 {
		logrus.Fatalf("数据库未配置")
	}
	dc := global.Config.DB[0]
	db, err := gorm.Open(postgres.Open(dc.PgsqlDsn()), &gorm.Config{
		Logger:                                   customGormLogger,
		DisableForeignKeyConstraintWhenMigrating: true, //不生成外建约束
	})
	if err != nil {
		logrus.Fatalf("数据库连接失败: %s", err)
	}
	sqlDB, err := db.DB()
	if sqlDB != nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
	logrus.Infof("数据库连接成功")

	if len(global.Config.DB) > 1 {
		var readList []gorm.Dialector
		for _, d := range global.Config.DB[1:] {
			readList = append(readList, mysql.Open(d.PgsqlDsn()))
		}
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(dc.PgsqlDsn())}, //读写
			Replicas: readList,                                    //读
			Policy:   dbresolver.RandomPolicy{},
		}))
		if err != nil {
			logrus.Fatalf("读写配置错误 %s", err)
		}
		logrus.Info("数据库读写配置成功")
	}
	if global.Config.DB[0].Debug && global.Config.System.Env == "dev" {
		return db.Debug()
	}
	return db
}
