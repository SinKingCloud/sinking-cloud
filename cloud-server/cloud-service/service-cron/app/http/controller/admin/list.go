package admin

import (
	"service-cron/app/repository/task"
	"service-cron/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// List 任务列表
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time,update_time", "id,create_time,update_time")
	var form struct {
		Name            string `json:"name" default:"" validate:"omitempty" label:"名称"`
		Type            string `json:"type" default:"" validate:"omitempty,numeric" label:"类型"`
		Status          string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
		Concurrent      string `json:"concurrent" default:"" validate:"omitempty,numeric" label:"允许并发"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &task.SelectTask{}
	if form.Name != "" {
		where.Name = &form.Name
	}
	if form.Type != "" {
		value, err := strconv.Atoi(form.Type)
		if err != nil {
			c.Error("类型参数错误")
			return
		}
		where.Type = &value
	}
	if form.Status != "" {
		value, err := strconv.Atoi(form.Status)
		if err != nil {
			c.Error("状态参数错误")
			return
		}
		where.Status = &value
	}
	if form.Concurrent != "" {
		value, err := strconv.Atoi(form.Concurrent)
		if err != nil {
			c.Error("并发配置参数错误")
			return
		}
		where.Concurrent = &value
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = &form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = &form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where.UpdateTimeStart = &form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where.UpdateTimeEnd = &form.UpdateTimeEnd
	}
	data, err := service.Task.Select(where, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", data)
}
