package main

import (
	models "gobee/models/base"
	"gobee/pkg/setting"

	_ "gobee/routers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	setting.Setup()
	models.Setup()
}

func main() {

	if web.BConfig.RunMode == "dev" {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	web.Run()
}
