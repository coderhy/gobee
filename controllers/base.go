package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	web.Controller
}

type ReturnMsg struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

func init() {

}

func (this *BaseController) ResponseJson(code int, msg string, data interface{}) {

	res := ReturnMsg{
		code, msg, data,
	}

	this.Data["json"] = res
	this.ServeJSON() //对json进行序列化输出
	this.StopRun()
}
