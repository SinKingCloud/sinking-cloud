package systems

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/config"
)

// List 获取菜单数据
func (ControllerConfig) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"sort" validate:"oneof=id sort update_time create_time" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Name            string `form:"name" default:"" validate:"max=20" label:"名称"`
		Describe        string `form:"describe" default:"" validate:"max=100" label:"描述"`
		Key             string `form:"key" default:"" validate:"max=50" label:"配置标识"`
		Value           string `form:"value" default:"" validate:"max=500" label:"配置内容"`
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
	if form.Key != "" {
		where["key"] = form.Key
	}
	if form.Value != "" {
		where["value"] = form.Value
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
	data, total, err := config.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取配置数据失败")
	} else {
		c.SuccessWithData("获取配置数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Update 修改信息
func (ControllerConfig) Update(c *sinking_server.Context) {
	type Form struct {
		Ids      []int  `json:"ids" default:"" validate:"required,min=1,unique" label:"ID列表"`
		Name     string `json:"name" default:"" validate:"omitempty,max=20" label:"名称"`
		Describe string `json:"describe" default:"" validate:"omitempty,max=100" label:"描述"`
		Value    string `json:"value" default:"" validate:"omitempty" label:"配置内容"`
		Sort     string `json:"sort" default:"" validate:"omitempty,numeric" label:"排序"`
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
	if form.Value != "" {
		data["value"] = form.Value
	}
	if form.Sort != "" {
		data["sort"] = form.Sort
	}
	err := config.GetIns().UpdateByIds(form.Ids, data)
	if err == nil {
		c.Success("修改数据成功")
	} else {
		c.Error("修改数据失败")
	}
}

// Create 添加信息
func (ControllerConfig) Create(c *sinking_server.Context) {
	type Form struct {
		Name     string `json:"name" validate:"required,max=20" label:"名称"`
		Describe string `json:"describe" validate:"required,max=100" label:"描述"`
		Key      string `json:"key" default:"" validate:"required,max=50" label:"配置标识"`
		Value    string `json:"value" default:"" validate:"required,max=500" label:"配置内容"`
		Sort     int    `json:"sort" default:"" validate:"omitempty,numeric" label:"排序"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	if num, err := config.GetIns().CountByKey(form.Key); err != nil || num > 0 {
		c.Error("该配置项名称已存在")
		return
	}
	err := config.GetIns().Create(&model.SystemConfig{
		Name:     form.Name,
		Describe: form.Describe,
		Key:      form.Key,
		Value:    form.Value,
		Sort:     form.Sort,
	})
	if err == nil {
		c.Success("添加数据成功")
	} else {
		c.Error("添加数据失败")
	}
}

// Delete 删除信息
func (ControllerConfig) Delete(c *sinking_server.Context) {
	type Form struct {
		Ids []int `json:"ids" default:"" validate:"required,min=1,unique" label:"ID列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := config.GetIns().DeleteByIds(form.Ids)
	if err == nil {
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}
