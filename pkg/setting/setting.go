package setting

import (
	"gobee/pkg/common"
	"gobee/pkg/utils"
	"io/ioutil"
	"path"
	"strings"

	"github.com/beego/beego/v2/core/config"
	logs "github.com/beego/beego/v2/core/logs"
)

func init() {
	utils.ConfigPath = config.DefaultString("configPath", "config") //yaml配置文件
}

// Setup initialize the configuration instance
func Setup() {

	if isExist := common.IsDirExist(utils.ConfigPath); isExist == false {
		logs.Error("配置目录路径:", utils.ConfigPath, "不存在,以加载默认配置启用:", utils.ConfigPath, "请检查配置")
		// fmt.Println("配置目录路径:", utils.ConfigPath, "不存在,以加载默认配置启用:", utils.ConfigPath, "请检查配置")
	}

	files, _ := ioutil.ReadDir(utils.ConfigPath)
	// var b bytes.Buffer
	for _, f := range files {
		var fileSuffix string
		fileSuffix = path.Ext(f.Name())
		var filenameOnly string
		filenameOnly = strings.TrimSuffix(f.Name(), fileSuffix)
		// utils.GetConfig("router")

		if utils.AllCacheConfig == nil {
			utils.AllCacheConfig = make(map[string]utils.ConfigEngine)
		}

		config := utils.ConfigEngine{}
		config.Load(utils.ConfigPath + "/" + f.Name())
		// return config
		//优化
		utils.AllCacheConfig[filenameOnly] = config
	}

}
