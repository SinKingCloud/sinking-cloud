package admin

import (
	"service-pay/app/repository/channel"
	"service-pay/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

type channelController struct {
}

// List 列表数据
func (channelController) List(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,create_time,update_time")
	var form struct {
		Name            string `json:"name" default:"" validate:"omitempty" label:"名称"`
		Plugin          string `json:"plugin" default:"" validate:"omitempty,numeric" label:"插件"`
		Status          string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	where := &channel.SelectChannel{}
	if form.Name != "" {
		where.Name = form.Name
	}
	if form.Plugin != "" {
		where.Plugin = form.Plugin
	}
	if form.Status != "" {
		where.Status = form.Status
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
	data, total, err := service.Channel.Select(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", page.NewPage(total, pageNum, pageSize, data))
}

// Info 详情数据
func (channelController) Info(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Channel.FindById(form.Id, false)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.SuccessWithData("获取数据成功", info)
}

// Update 更新数据
func (channelController) Update(c *context.Context) {
	var form struct {
		Id     int64  `json:"id" default:"" validate:"required,min=1" label:"记录ID"`
		Name   string `json:"name" default:"" validate:"omitempty" label:"名称"`
		Config string `json:"config" default:"" validate:"omitempty" label:"配置"`
		Status string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &channel.UpdateChannel{}
	if form.Name != "" {
		data.Name = form.Name
	}
	if form.Config != "" {
		data.Config = form.Config
	}
	if form.Status != "" {
		data.Status, _ = strconv.Atoi(form.Status)
	}
	err := service.Channel.UpdateById(form.Id, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

// Create 创建数据
func (channelController) Create(c *context.Context) {
	var form struct {
		Name   string `json:"name" default:"" validate:"required" label:"通道名称"`
		Type   string `json:"type" default:"" validate:"required" label:"通道类型"`
		Plugin int    `json:"plugin" default:"" label:"使用插件"`
		Status int    `json:"status" default:"" label:"通道状态"`
		Config string `json:"config" default:"" label:"通道配置"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Channel.Create(form.Name, form.Type, form.Plugin, form.Status, form.Config)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("创建数据成功")
}

// Delete 删除数据
func (channelController) Delete(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Channel.DeleteById(form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("删除数据成功")
}
