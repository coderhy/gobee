package models

import (
	"github.com/astaxie/beego/orm"
)

func Setup() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	maxIdle := 30
	maxConn := 30
	orm.RegisterDataBase("default", "mysql", "root:`c2#^@j1T)oX:_@tcp(w.talk-test.int.yidian-inc.com)/shenbian_ugc?charset=utf8", orm.MaxIdleConnections(maxIdle), orm.MaxOpenConnections(maxConn))
}
