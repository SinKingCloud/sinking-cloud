package web

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/web_notice"
)

// List 获取数据
func (ControllerNotice) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"sort" validate:"oneof=id sort create_time update_time" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Title           string `form:"title" default:"" validate:"omitempty" label:"标题"`
		Status          string `form:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
		Place           string `form:"place" default:"" validate:"omitempty" label:"公告类型"`
		CreateTimeStart string `form:"create_time_start" default:"" validate:"-" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" default:"" validate:"-" label:"创建结束时间"`
		UpdateTimeStart string `form:"update_time_start" default:"" validate:"-" label:"更新起始时间"`
		UpdateTimeEnd   string `form:"update_time_end" default:"" validate:"-" label:"更新结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.Place != "" {
		where["place"] = form.Place
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where["update_time_start"] = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where["update_time_end"] = form.UpdateTimeEnd
	}
	if form.Title != "" {
		where["title"] = form.Title
	}
	data, total, err := web_notice.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取公告列表失败")
	} else {
		system_log.GetIns().CreateByQueue(c, system_log.EventShow, "查看公告", "查看公告列表")
		c.SuccessWithData("获取公告列表成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Info 获取详情
func (ControllerNotice) Info(c *sinking_server.Context) {
	type Form struct {
		Id int `form:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	data, err := web_notice.GetIns().FindById(form.Id)
	if err == nil && data != nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventShow, "查看公告", "查看公告"+data.Title)
		c.SuccessWithData("获取详情数据成功", data)
	} else {
		c.Error("获取详情数据失败")
	}
}

// Update 修改信息
func (ControllerNotice) Update(c *sinking_server.Context) {
	type Form struct {
		Ids     []int  `json:"ids" default:"" validate:"required,min=1,unique" label:"ID列表"`
		Place   string `json:"place" default:"" validate:"omitempty,max=20" label:"位置"`
		Title   string `json:"title" default:"" validate:"omitempty,max=100" label:"名称"`
		Content string `json:"content" default:"" validate:"omitempty" label:"内容"`
		Sort    string `json:"sort" default:"" validate:"omitempty,numeric" label:"排序"`
		Status  string `json:"status" default:"" validate:"omitempty,numeric,oneof=0 1" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	data := make(map[string]interface{})
	if form.Place != "" {
		data["place"] = form.Place
	}
	if form.Title != "" {
		data["title"] = form.Title
	}
	if form.Content != "" {
		data["content"] = form.Content
	}
	if form.Sort != "" {
		data["sort"] = form.Sort
	}
	if form.Status != "" {
		data["status"] = form.Status
	}
	err := web_notice.GetIns().UpdateByIds(form.Ids, data)
	if err == nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventEdit, "修改公告", "修改公告信息")
		c.Success("修改数据成功")
	} else {
		c.Error("修改数据失败")
	}
}

// Create 添加信息
func (ControllerNotice) Create(c *sinking_server.Context) {
	type Form struct {
		Place   string `json:"place" default:"" validate:"required,max=20" label:"位置"`
		Title   string `json:"title" default:"" validate:"required,max=100" label:"名称"`
		Content string `json:"content" default:"" validate:"required" label:"内容"`
		Sort    int    `json:"sort" default:"" validate:"numeric" label:"排序"`
		Status  int    `json:"status" validate:"oneof=0 1" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := web_notice.GetIns().Create(&model.WebNotice{
		Place:   form.Place,
		Title:   form.Title,
		Content: form.Content,
		Status:  form.Status,
		Sort:    form.Sort,
		LookNum: 0,
	})
	if err == nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventCreate, "添加公告", "添加公告信息")
		c.Success("添加数据成功")
	} else {
		c.Error("添加数据失败")
	}
}

// Delete 删除信息
func (ControllerNotice) Delete(c *sinking_server.Context) {
	type Form struct {
		Ids []int `json:"ids" default:"" validate:"required,min=1,unique" label:"ID列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := web_notice.GetIns().DeleteByIds(form.Ids)
	if err == nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventDelete, "删除公告", "删除公告数据")
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}
