package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
)

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(TopicModel))
}

type TopicModel struct {
	Id          int
	Name        string `orm:"size(60)" form:"Name"  valid:"Required;MaxSize(60);MinSize(6)"`
	Icon        string `orm:"size(100)" form:"Icon" valid:"Required;MaxSize(100);MinSize(6)"`
	Content     string `orm:"size(100)" form:"Content" valid:"Required"`
	Status      int    `orm:"default(1)" form:"Status" valid:"Required;MaxSize(1);MinSize(0)"`
	Weight      int    `orm:"size(11);default(0)" form:"Weight"`
	CreateTime  int    `orm:"size(11)" form:"CreateTime"`
	PublishTime int    `orm:"size(11)" form:"PublishTime"`
}

func (t *TopicModel) TableName() string {
	return "topic"
}

func GetOne(topicID int) (t *TopicModel, err error) {
	o := orm.NewOrm()
	topic := TopicModel{Id: topicID}

	err = o.Read(&topic)

	if err == orm.ErrNoRows {
		return &TopicModel{}, errors.New("查询不到")
	} else if err == orm.ErrMissPK {
		return &TopicModel{}, errors.New("找不到主键")
	}
	return &topic, nil
}
