package topicservice

import (
	"context"
	models "gobee/models/topic"
	"gobee/pkg/common"
	"time"
)

type Topic struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (t *Topic) GetTopic() (*models.TopicModel, error) {
	topic, err := models.GetOne(t.ID)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func hardWork(job interface{}) error {
	time.Sleep(time.Second * 4)
	return nil
}

func RequestWork(ctx context.Context, job interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- hardWork(job)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func hardWork2(job interface{}) error {
	panic("hy")
}

//recover只能捕获当前协程当前函数或直接调用函数的panic  其他协程panic无法直接捕获
func RequestWork2(ctx context.Context, job interface{}) error {
	// ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	ctx, cancel := common.ShrinkDeadline(ctx, time.Second*2)
	defer cancel()

	done := make(chan error, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		done <- hardWork2(job)
	}()

	select {
	case err := <-done:
		return err
	case p := <-panicChan:
		panic(p)
	case <-ctx.Done():
		return ctx.Err()
	}
}
