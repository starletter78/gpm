package core

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gpm/conf"
	"gpm/flags"
	"gpm/global"
)

func ReadConf() (c *conf.Config) {
	viper.SetConfigName(flags.FlagOptions.File)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("conf")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal(err)
	}
	if err := viper.Unmarshal(&global.Config); err != nil {
		logrus.Fatal("反序列化配置失败: ", err)
	}
	return
}
