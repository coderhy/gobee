package models

import (
	"fmt"
	"gobee/pkg/utils"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

func Setup() {
	ConfigEngine := utils.AllCacheConfig["mysql"]
	mysqlData := ConfigEngine.GetAll()
	// fmt.Println(reflect.ValueOf(mysqlData).Kind())

	orm.RegisterDriver("mysql", orm.DRMySQL)
	for k, v := range mysqlData {
		dbAlias := k.(string)
		data := v.(map[interface{}]interface{})
		maxIdle := data["MaxIdleConns"].(int)
		maxConn := data["MaxIdleConns"].(int)
		dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", data["DbUser"], data["DbPwd"], data["DbHost"], data["DbPort"], data["DbName"], data["DbCharset"])

		fmt.Println(dbDSN)
		orm.RegisterDataBase(dbAlias, "mysql", dbDSN, orm.MaxIdleConnections(maxIdle), orm.MaxOpenConnections(maxConn))
	}

}
