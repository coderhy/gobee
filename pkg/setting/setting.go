package setting

import (
	_ "github.com/beego/beego/v2/core/config/toml"
)

var (
	//ConfigPath 全局配置路径
	ConfigPath string
)

// Setup initialize the configuration instance
func Setup() {
	// ConfigPath := beego.AppConfig.DefaultString("ConfigPath2", "conf")

	// if isExist := common.IsDirExist(ConfigPath); isExist == false {
	// 	fmt.Println("配置目录路径:", ConfigPath, "不存在,以加载默认配置启用:", ConfigPath, "请检查配置")
	// }

	// files, _ := ioutil.ReadDir(ConfigPath)
	// var b bytes.Buffer
	// for _, f := range files {
	// 	fmt.Println(f.Name())
	// if ok := beego.LoadAppConfig("ini", ConfigPath+"/"+f.Name()); ok != nil {
	// 	log.Println("配置文件:", f.Name(), "不存在,请检查配置文件")
	// }
	// if ok := config.InitGlobalInstance("toml", ConfigPath+"/"+f.Name()); ok != nil {
	// 	log.Println("配置文件:", f.Name(), "不存在,请检查配置文件")
	// }
	// ctx2, _ := ioutil.ReadFile(ConfigPath + "/" + f.Name())

	// fmt.Println(ctx2)
	// ww, ok := b.Write(ctx2)
	// }

	// ctx2, _ := ioutil.ReadFile(ConfigPath + "/mysql.toml")
	// config.NewConfigData("toml", ctx2)

	// val, err := config.GetSection("Mysql")
	// log.Println(val, err)
	// val2, err := config.GetSection("Mongo.ShenbianUgc")
	// log.Println(val2, err)

	// val3, err := beego.AppConfig.GetSection("mongo.shenbian_ugc")
	// log.Println(val3, err)
	// val4, err := beego.AppConfig.GetSection("mysql.shenbian_ugc")
	// log.Println(val4, err)

	// val5, err := beego.AppConfig.GetSection("shenbianugc")
	// log.Println(val5, err)

	// val6, err := beego.AppConfig.GetSection("mongo.shenbian_ugc")
	// log.Println(val6, err)
	// val5, err := beego.AppConfig.GetSection("mysql.shenbian_ugc")
	// log.Println(val5, err)

}
