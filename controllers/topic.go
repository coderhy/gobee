package controllers

import (
	"encoding/json"
	"fmt"
	"gobee/pkg/utils"
	"log"

	"github.com/astaxie/beego/validation"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type TopicController struct {
	beego.Controller
}

type Topic struct {
	Name string `form:"name"; valid:"Required;MaxSize(60)"`
	Icon string `form:"icon"; valid:"Required;MaxSize(100)"`
}

type GetOneRule struct {
	TopicID int `valid:"Required;Min(1)"`
}

type Userinfo struct {
	DbHost       string `json:"dbhost"`
	DbUser       string `json:"dbuser"`
	DbPwd        string `json:"dbpwd"`
	DbName       string `json:"dbname"`
	DbPort       string `json:"dbport"`
	MaxOpenConns int    `json:"maxopenconns"`
	MaxIdleConns int    `json:"maxidleconns"`
	Debug        bool   `json:"debug"`
}

// @Title GetTopic
// @Description find topic by topicID
// @Param	body		body 	GetOneRule	true		"The rule"
// @Success 200 {object} GetOneRule
// @Failure 403 body is empty
// @router /GetTopic [post]
func (t *TopicController) GetTopic() {

	var form GetOneRule
	json.Unmarshal(t.Ctx.Input.RequestBody, &form)
	log.Println(form)
	valid := validation.Validation{}
	_, err := valid.Valid(&form)
	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			t.Data["json"] = err.Message
		}
	}
	if err != nil {
		t.Data["json"] = err.Error()
	}

	mongoData := utils.AllCacheConfig["mongo"]
	demoData := mongoData.Get("MG_ALL_DEMO").(map[interface{}]interface{})

	result := map[string]interface{}{}
	data := map[string]interface{}{}
	result["code"] = 1
	result["msg"] = "成功"

	for k, v := range demoData {
		data[k.(string)] = v
	}

	aa, _ := json.Marshal(data)
	fmt.Println(string(aa))
	var userinfo Userinfo
	json.Unmarshal(aa, &userinfo)
	result["data"] = userinfo
	t.Data["json"] = result

	// if !b {
	// 	topicInfo, err := models.GetOne(form.TopicID)
	// 	if err != nil {
	// 		t.Data["json"] = err.Error()
	// 	} else {
	// 		t.Data["json"] = topicInfo
	// 	}
	// }
	t.ServeJSON()
}

//Prepare @Title  预
func (u *TopicController) Prepare() {
	l := logs.GetLogger()
	l.Println("我先走一步")
	// u.Data["json"] = "Prepare"
	// u.ServeJSON()
}
