package scene

import (
	"service-driver/app/repository/scene"
	"service-driver/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// List 列表数据
func (Controller) List(c *context.Context) {
	pageNum, pageSize, orderByField, orderByType := c.ValidatePage("id", "desc", "id,create_time,update_time")
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
		where.Name = form.Name
	}
	if form.ChannelId != "" {
		where.ChannelId = form.ChannelId
	}
	if form.Login != "" {
		where.Login = form.Login
	}
	if form.Status != "" {
		where.Status = form.Status
	}
	if form.Remark != "" {
		where.Remark = form.Remark
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where.UpdateTimeStart = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where.UpdateTimeEnd = form.UpdateTimeEnd
	}
	data, total, err := service.Scene.Select(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}
