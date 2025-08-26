package conf

import (
	"fmt"
)

type DB struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DB       string `yaml:"db"`
	Debug    bool   `yaml:"debug"`  //打印全部日志
	Source   string `yaml:"source"` //源
}

func (d DB) MysqlDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DB)
}

func (d DB) PgsqlDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", d.Host, d.User, d.Password, d.DB)
}

func (d DB) Empty() bool {
	return d.User == "" && d.Password == "" && d.Host == "" && d.Port == 0
}
func (d DB) Addr() string {
	return fmt.Sprintf("%s:%d", d.Host, d.Port)
}
