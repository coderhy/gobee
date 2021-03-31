package models

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(TopicModel))
}

type TopicModel struct {
	Id          int    `json:"id" orm:"size(11);column(id)"`
	Name        string `json:"name" orm:"size(60);column(name)"`
	Icon        string `json:"icon" orm:"size(100);column(icon)"`
	Content     string `json:"content" orm:"size(100);column(content)"`
	Status      int    `json:"status" orm:"default(1);column(status)"`
	Weight      int    `json:"weight" orm:"size(11);column(weight);default(0)"`
	CreateTime  int    `json:"createTime" orm:"size(11);column(create_time)"`
	PublishTime int    `json:"publishTime" orm:"size(11);column(publish_time)"`
}

func (t *TopicModel) TableName() string {
	return "topic"
}

func (t *TopicModel) GetDBName() string {
	return "DB_UGC"
}

func GetOne(topicID int) (t *TopicModel, err error) {
	// o := orm.NewOrm()
	o := orm.NewOrmUsingDB(t.GetDBName())
	topic := TopicModel{Id: topicID}

	err = o.Read(&topic)

	if err == orm.ErrNoRows {
		return &TopicModel{}, errors.New("查询不到")
	} else if err == orm.ErrMissPK {
		return &TopicModel{}, errors.New("找不到主键")
	}
	return &topic, nil
}
