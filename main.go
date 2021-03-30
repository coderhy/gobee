package main

import (
	"gobee/pkg/setting"
	_ "gobee/routers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	setting.Setup()
}

func main() {

	if web.BConfig.RunMode == "dev" {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	web.Run()
}
