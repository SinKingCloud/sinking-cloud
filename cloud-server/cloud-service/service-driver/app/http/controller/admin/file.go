package admin

import (
	"service-driver/app/repository/file"
	"service-driver/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

type fileController struct {
}

// List 列表数据
func (fileController) List(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,create_time,update_time")
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
	data, total, err := service.File.Select(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}

// Update 更新数据
func (fileController) Update(c *context.Context) {
	var form struct {
		Ids    []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
		Status string  `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &file.UpdateFile{}
	if form.Status != "" {
		data.Status, _ = strconv.Atoi(form.Status)
	}
	err := service.File.UpdateByIds(form.Ids, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

// Delete 删除数据
func (fileController) Delete(c *context.Context) {

}
