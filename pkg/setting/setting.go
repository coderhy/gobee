package setting

import (
	"fmt"
	"gobee/pkg/common"
	"gobee/pkg/utils"
	"io/ioutil"
	"path"
	"strings"

	"github.com/beego/beego/v2/core/config"
	// beego "github.com/beego/beego/v2/server/web"
)

// Setup initialize the configuration instance
func Setup() {
	utils.ConfigPath = config.DefaultString("configPath", "config")
	if isExist := common.IsDirExist(utils.ConfigPath); isExist == false {
		fmt.Println("配置目录路径:", utils.ConfigPath, "不存在,以加载默认配置启用:", utils.ConfigPath, "请检查配置")
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
