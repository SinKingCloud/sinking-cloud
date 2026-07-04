package file

import (
	"service-driver/app/repository/file"
	"service-driver/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// List 列表数据
func (Controller) List(c *context.Context) {
	query := c.ValidatePage("id", "desc", "id,create_time,update_time")
	var form struct {
		ChannelId       string `json:"channel_id" default:"" validate:"omitempty,numeric" label:"通道ID"`
		SceneId         string `json:"scene_id" default:"" validate:"omitempty,numeric" label:"场景ID"`
		Ext             string `json:"ext" default:"" validate:"omitempty" label:"类型"`
		Hash            string `json:"hash" default:"" validate:"omitempty" label:"哈希"`
		Status          string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
		Name            string `json:"name" default:"" validate:"omitempty" label:"名称"`
		Path            string `json:"path" default:"" validate:"omitempty" label:"路径"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &file.SelectFile{}
	if form.ChannelId != "" {
		where.ChannelId = form.ChannelId
	}
	if form.SceneId != "" {
		where.SceneId = form.SceneId
	}
	if form.Ext != "" {
		where.Ext = form.Ext
	}
	if form.Hash != "" {
		where.Hash = form.Hash
	}
	if form.Status != "" {
		where.Status = form.Status
	}
	if form.Name != "" {
		where.Name = form.Name
	}
	if form.Path != "" {
		where.Path = form.Path
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
	data, total, err := service.File.Select(where, query)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, query.Page, query.PageSize, data))
}
