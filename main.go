package main

import (
	"gpm/app/router"
	"gpm/core"
	"gpm/flags"
	"gpm/global"
)

func main() {
	flags.Parse()
	core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()
	flags.Run()
	global.CasbinEnforcer = core.InitCasbin()
	router.Run()
}
