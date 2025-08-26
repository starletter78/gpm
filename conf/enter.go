package conf

type Config struct {
	System    System `yaml:"system"`
	Log       Log    `yaml:"log"`
	DB        []DB   `yaml:"db"` //数据库连接列表
	Jwt       Jwt    `yaml:"jwt"`
	ArgsCheck ArgsCheck
}
