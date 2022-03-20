package systems

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/access"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/log"
)

// List 接口列表
func (ControllerAccess) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id update_time create_time status" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Name            string `form:"name" default:"" validate:"max=20" label:"名称"`
		Describe        string `form:"describe" default:"" validate:"max=100" label:"描述"`
		Url             string `form:"url" default:"" validate:"max=50" label:"请求地址"`
		Method          string `form:"method" default:"" validate:"omitempty,oneof=GET POST PUT UPDATE" label:"请求类型"`
		Status          string `form:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
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
	if form.Name != "" {
		where["name"] = form.Name
	}
	if form.Describe != "" {
		where["describe"] = form.Describe
	}
	if form.Url != "" {
		where["url"] = form.Url
	}
	if form.Method != "" {
		where["method"] = form.Method
	}
	if form.Status != "" {
		where["status"] = form.Status
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
	data, total, err := access.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err == nil {
		log.GetIns().CreateByQueue(c, 0, "查看权限", "查看权限列表")
		c.SuccessWithData("获取权限数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
		return
	}
	c.Error("获取数据失败")
}

// Update 修改接口信息
func (ControllerAccess) Update(c *sinking_server.Context) {
	type Form struct {
		Ids      []int  `json:"ids" default:"" validate:"required,min=1,unique" label:"ID列表"`
		Name     string `json:"name" default:"" validate:"omitempty,max=20" label:"名称"`
		Describe string `json:"describe" default:"" validate:"omitempty,max=100" label:"描述"`
		Url      string `json:"url" default:"" validate:"omitempty,max=50" label:"请求地址"`
		Method   string `json:"method" default:"" validate:"omitempty,oneof=GET POST PUT UPDATE" label:"请求类型"`
		Status   string `json:"status" default:"" validate:"omitempty,numeric,oneof=0 1" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	data := make(map[string]interface{})
	if form.Name != "" {
		data["name"] = form.Name
	}
	if form.Describe != "" {
		data["describe"] = form.Describe
	}
	if form.Url != "" {
		data["url"] = form.Url
	}
	if form.Method != "" {
		data["method"] = form.Method
	}
	if form.Status != "" {
		data["status"] = form.Status
	}
	err := access.GetIns().UpdateByIds(form.Ids, data)
	if err == nil {
		log.GetIns().CreateByQueue(c, 3, "修改权限", "修改权限信息")
		c.Success("修改数据成功")
	} else {
		c.Error("修改数据失败")
	}
}

// Create 添加接口信息
func (ControllerAccess) Create(c *sinking_server.Context) {
	type Form struct {
		Name     string `json:"name" validate:"required,max=20" label:"名称"`
		Describe string `json:"describe" validate:"required,max=100" label:"描述"`
		Url      string `json:"url" validate:"required,max=50" label:"请求地址"`
		Method   string `json:"method" validate:"required,oneof=GET POST PUT UPDATE" label:"请求类型"`
		Status   int    `json:"status" validate:"oneof=0 1" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := access.GetIns().Create(&model.SystemAccess{
		Name:     form.Name,
		Describe: form.Describe,
		Url:      form.Url,
		Method:   form.Method,
		Status:   form.Status,
	})
	if err == nil {
		log.GetIns().CreateByQueue(c, 1, "添加权限", "添加权限信息")
		c.Success("添加数据成功")
	} else {
		c.Error("添加数据失败")
	}
}

// Delete 删除接口信息
func (ControllerAccess) Delete(c *sinking_server.Context) {
	type Form struct {
		Ids []int `json:"ids" default:"" validate:"required,min=1,unique" label:"ID列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := access.GetIns().DeleteByIds(form.Ids)
	if err == nil {
		log.GetIns().CreateByQueue(c, 2, "删除权限", "删除权限信息")
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}
