package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"

	base "gobee/controllers/v1/base"
	"gobee/pkg/e"
	topicservice "gobee/service/topic-service"
	"log"
	"sync"
	"time"

	// "github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

// Operations about Users
type TopicController struct {
	// beego.Controller
	base.BaseController
}

type GetOneRule struct {
	TopicID int `valid:"Required;Min(1)"`
}

type TopicInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("GetTopic: %v", err)
			t.ResponseJson(e.ERROR, err.(string), map[string]interface{}{})
		}
	}()

	var form GetOneRule
	json.Unmarshal(t.Ctx.Input.RequestBody, &form)
	valid := validation.Validation{}
	check, err := valid.Valid(&form)
	if err != nil {
		t.Data["json"] = err.Error()
	}
	if !check {
		// 如果有错误信息，证明验证没通过
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)

			//1.t.ServeJSON()
			t.Data["json"] = map[string]interface{}{
				"code":   e.ERROR,
				"msg":    e.GetMsg(e.ERROR),
				"result": map[string]interface{}{},
			}
			t.ServeJSON() //对json进行序列化输出

			// 2.recover panic
			// panic(err.Message)
		}
	}

	//获取topic详情
	topicService := topicservice.Topic{ID: form.TopicID}
	topicInfo, err := topicService.GetTopic()
	if err != nil {
		t.ResponseJson(e.ERROR, err.Error(), map[string]interface{}{})
	}
	t.ResponseJson(e.SUCCESS, "success", topicInfo)
}

// @Title GetTopicAll  context超时处理 demo hhh
// @Description find topic by topicID
// @Param	body		body 	GetOneRule	true		"The rule"
// @Success 200 {object} GetOneRule
// @Failure 403 body is empty
// @router /GetTopicAll [get]
func (t *TopicController) GetTopicAll() {
	const total = 1000
	var wg sync.WaitGroup
	wg.Add(total)
	now := time.Now()
	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			topicservice.RequestWork(context.Background(), "any")
		}()
	}
	wg.Wait()
	fmt.Println("elapsed:", time.Since(now))
	time.Sleep(time.Second * 8)
	fmt.Println("number of goroutines:", runtime.NumGoroutine())
}

// @Title GetTopicPanic  捕获panic
// @Description find topic by topicID
// @Param	body		body 	GetOneRule	true		"The rule"
// @Success 200 {object} GetOneRule
// @Failure 403 body is empty
// @router /GetTopicPanic [get]
func (t *TopicController) GetTopicPanic() {
	const total = 10
	var wg sync.WaitGroup
	wg.Add(total)
	now := time.Now()
	for i := 0; i < total; i++ {
		go func() {
			defer func() {
				if p := recover(); p != nil {
					fmt.Println("oops, panic")
				}
			}()

			defer wg.Done()
			topicservice.RequestWork2(context.Background(), "any")
		}()
	}
	wg.Wait()
	fmt.Println("elapsed:", time.Since(now))
	time.Sleep(time.Second * 20)
	fmt.Println("number of goroutines:", runtime.NumGoroutine())
}

//Prepare @Title  预
func (u *TopicController) Prepare() {

	// l := logs.GetLogger()
	// l.Println("我先走一步")
	// u.Data["json"] = "Prepare"
	// u.ServeJSON()
}
