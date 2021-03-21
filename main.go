package main

import (
	"gobee/pkg/setting"
	_ "gobee/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	setting.Setup()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
