package admin

import (
	"service-driver/app/model"
	"service-driver/app/repository/scene"
	"service-driver/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

type sceneController struct {
}

// List 列表数据
func (sceneController) List(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,create_time,update_time")
	var form struct {
		Name            string `json:"name" default:"" validate:"omitempty" label:"名称"`
		ChannelId       string `json:"channel_id" default:"" validate:"omitempty,numeric" label:"通道ID"`
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

// Update 更新数据
func (sceneController) Update(c *context.Context) {
	var form struct {
		Id       int64  `json:"id" default:"" validate:"required,min=1" label:"记录ID"`
		Name     string `json:"name" default:"" validate:"omitempty" label:"名称"`
		MaxSize  string `json:"max_size" default:"" validate:"omitempty,numeric" label:"最大限制"`
		AllowExt string `json:"allow_ext" default:"" validate:"omitempty" label:"允许类型"`
		Remark   string `json:"remark" default:"" validate:"omitempty" label:"备注"`
		Status   string `json:"status" default:"" validate:"omitempty,numeric" label:"状态"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &scene.UpdateScene{}
	if form.Name != "" {
		data.Name = form.Name
	}
	if form.MaxSize != "" {
		data.MaxSize, _ = strconv.ParseInt(form.MaxSize, 10, 64)
	}
	if form.AllowExt != "" {
		data.AllowExt = form.AllowExt
	}
	if form.Remark != "" {
		data.Remark = form.Remark
	}
	if form.Status != "" {
		data.Status, _ = strconv.Atoi(form.Status)
	}
	err := service.Scene.UpdateById(form.Id, data)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

// Create 创建数据
func (sceneController) Create(c *context.Context) {
	var form struct {
		Name      string `json:"name" default:"" validate:"required" label:"名称"`
		MaxSize   int64  `json:"max_size" default:"0" validate:"required" label:"最大限制"`
		AllowExt  string `json:"allow_ext" default:"*" validate:"omitempty" label:"允许类型"`
		ChannelId int64  `json:"channel_id" default:"" validate:"required" label:"通道ID"`
		Status    int    `json:"status" default:"0" validate:"required" label:"状态"`
		Remark    string `json:"remark" default:"" validate:"omitempty" label:"备注"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Scene.Create(&model.Scene{
		Name:      form.Name,
		MaxSize:   form.MaxSize,
		AllowExt:  form.AllowExt,
		ChannelId: form.ChannelId,
		Remark:    form.Remark,
		Status:    form.Status,
	})
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("创建数据成功")
}

// Delete 删除数据
func (sceneController) Delete(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Scene.DeleteById(form.Id)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("删除数据成功")
}
