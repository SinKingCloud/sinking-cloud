package admin

import (
	"fmt"
	"service-cron/app/repository/task"
	"service-cron/app/service"
	"strconv"

	commonService "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
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
		data.Name = &form.Name
	}
	if form.Type != "" {
		value, err := strconv.Atoi(form.Type)
		if err != nil {
			c.Error("类型参数错误")
			return
		}
		data.Type = &value
	}
	if form.Spec != "" {
		data.Spec = &form.Spec
	}
	if form.Script != "" {
		data.Script = &form.Script
	}
	if form.Status != "" {
		value, err := strconv.Atoi(form.Status)
		if err != nil {
			c.Error("状态参数错误")
			return
		}
		data.Status = &value
	}
	if form.Concurrent != "" {
		value, err := strconv.Atoi(form.Concurrent)
		if err != nil {
			c.Error("并发配置参数错误")
			return
		}
		data.Concurrent = &value
	}
	err := service.Task.UpdateById(form.Id, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = commonService.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventEdit, "修改计划任务", fmt.Sprintf("修改计划任务[%d]", form.Id))
	c.Success("更新数据成功")
}
