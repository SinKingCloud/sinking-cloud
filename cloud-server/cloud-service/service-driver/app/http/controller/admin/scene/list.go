package scene

import (
	"service-driver/app/repository/scene"
	"service-driver/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// List 列表数据
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time,update_time", "id,create_time,update_time")
	var form struct {
		Name            string `json:"name" default:"" validate:"omitempty" label:"名称"`
		ChannelId       string `json:"channel_id" default:"" validate:"omitempty,numeric" label:"通道ID"`
		Login           string `json:"login" default:"" validate:"omitempty,numeric" label:"登录限制"`
		Status          string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
		Remark          string `json:"remark" default:"" validate:"omitempty" label:"名称"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &scene.SelectScene{}
	if form.Name != "" {
		where.Name = &form.Name
	}
	if form.ChannelId != "" {
		channelId, err := strconv.ParseInt(form.ChannelId, 10, 64)
		if err != nil {
			c.Error("通道ID参数错误")
			return
		}
		where.ChannelId = &channelId
	}
	if form.Login != "" {
		login, err := strconv.Atoi(form.Login)
		if err != nil {
			c.Error("登录限制参数错误")
			return
		}
		where.Login = &login
	}
	if form.Status != "" {
		status, err := strconv.Atoi(form.Status)
		if err != nil {
			c.Error("状态参数错误")
			return
		}
		where.Status = &status
	}
	if form.Remark != "" {
		where.Remark = &form.Remark
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
	data, err := service.Scene.Select(where, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", data)
}
