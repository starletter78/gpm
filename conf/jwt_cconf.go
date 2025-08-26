package conf

type Jwt struct {
	AccessExpire  int    `yaml:"accessExpire"`
	RefreshExpire int    `yaml:"refreshExpire"`
	AccessSecret  string `yaml:"accessSecret"`
	RefreshSecret string `yaml:"refreshSecret"`
	Issuer        string `yaml:"issuer"`
}
