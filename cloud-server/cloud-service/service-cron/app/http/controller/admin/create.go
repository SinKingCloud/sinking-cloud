package admin

import (
	"service-cron/app/model"
	"service-cron/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Create 创建任务
func (Controller) Create(c *context.Context) {
	var form struct {
		Name       string `json:"name" default:"" validate:"required" label:"任务名称"`
		Type       int    `json:"type" default:"0" validate:"required" label:"任务类型"`
		Spec       string `json:"spec" default:"" validate:"required" label:"任务表达式"`
		Script     string `json:"script" default:"" validate:"required" label:"任务内容"`
		Status     int    `json:"status" default:"0" validate:"required" label:"任务状态"`
		Concurrent int    `json:"concurrent" default:"0" validate:"required" label:"允许并发"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Task.Create(&model.Task{
		Type:       form.Type,
		Name:       form.Name,
		Spec:       form.Spec,
		Script:     form.Spec,
		Concurrent: form.Concurrent,
		Status:     form.Status,
	})
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("创建数据成功")
}
