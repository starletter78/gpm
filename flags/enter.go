package flags

import (
	"flag"
	"os"
)

type Options struct {
	File    string
	DB      bool
	Version bool
}

var FlagOptions = new(Options)

func Parse() {
	flag.BoolVar(&FlagOptions.DB, "db", false, "数据库迁移")
	flag.StringVar(&FlagOptions.File, "f", "settings.yaml", "配置文件")
	flag.BoolVar(&FlagOptions.Version, "v", false, "版本")
	flag.Parse()
}
func Run() {
	if FlagOptions.DB {
		FlagsDb()
		os.Exit(0)
	}
}
