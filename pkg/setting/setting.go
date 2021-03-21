package setting

import (
	"flag"
	"log"

	"github.com/beego/beego/v2/core/config"
	_ "github.com/beego/beego/v2/core/config/toml"
)

var (
	//ConfigPath 全局配置路径
	ConfigPath string
)

// Setup initialize the configuration instance
func Setup() {
	flag.StringVar(&ConfigPath, "c", "config/", "默认配置目录")
	defaultConfigPath := ConfigPath

	err := config.InitGlobalInstance("toml", defaultConfigPath+"/mysql.toml")
	log.Println(err)
	val, err := config.GetSection("database")
	log.Println(val, err)
}
