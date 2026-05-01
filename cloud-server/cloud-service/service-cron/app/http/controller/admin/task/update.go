package task

import (
	"service-cron/app/repository/task"
	"service-cron/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Update 修改任务
func (Controller) Update(c *context.Context) {
	var form struct {
		Id         int64  `json:"id" default:"" validate:"required,min=1" label:"任务ID"`
		Name       string `json:"name" default:"" validate:"omitempty" label:"任务名称"`
		Type       string `json:"type" default:"" validate:"omitempty,numeric" label:"任务类型"`
		Spec       string `json:"spec" default:"" validate:"omitempty" label:"任务表达式"`
		Script     string `json:"script" default:"" validate:"omitempty" label:"任务内容"`
		Status     string `json:"status" default:"" validate:"omitempty,numeric" label:"任务状态"`
		Concurrent string `json:"concurrent" default:"" validate:"omitempty,numeric" label:"允许并发"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &task.UpdateTask{}
	if form.Name != "" {
		data.Name = form.Name
	}
	if form.Type != "" {
		data.Type, _ = strconv.Atoi(form.Type)
	}
	if form.Spec != "" {
		data.Spec = form.Spec
	}
	if form.Script != "" {
		data.Script = form.Script
	}
	if form.Status != "" {
		data.Status, _ = strconv.Atoi(form.Status)
	}
	if form.Concurrent != "" {
		data.Concurrent, _ = strconv.Atoi(form.Concurrent)
	}
	err := service.Task.UpdateById(form.Id, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}
