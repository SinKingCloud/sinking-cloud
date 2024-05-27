package systems

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
)

// List 获取菜单数据
func (ControllerConfig) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"sort" validate:"oneof=id sort update_time create_time" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Pid             string `form:"pid" json:"pid" default:"" validate:"omitempty,numeric" label:"父级ID"`
		Name            string `form:"name" json:"name" default:"" validate:"omitempty,max=20" label:"名称"`
		Describe        string `form:"describe" json:"describe" default:"" validate:"omitempty,max=100" label:"描述"`
		Key             string `form:"key" json:"key" default:"" validate:"omitempty,max=50" label:"配置标识"`
		Value           string `form:"value" json:"value" default:"" validate:"omitempty,max=500" label:"配置内容"`
		Control         string `form:"control" json:"control" default:"" validate:"omitempty,max=500" label:"控件类型"`
		CreateTimeStart string `form:"create_time_start" json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `form:"create_time_end" json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `form:"update_time_start" json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `form:"update_time_end" json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.Pid != "" {
		where["pid"] = form.Pid
	}
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
	if form.Control != "" {
		where["control"] = form.Control
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
	data, total, err := service2.SystemConfig.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取配置数据失败")
	} else {
		service2.SystemLog.CreateByQueue(c, system_log.EventShow, "查看配置", "查看配置列表")
		c.SuccessWithData("获取配置数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Update 修改信息
func (ControllerConfig) Update(c *sinking_server.Context) {
	type Form struct {
		Ids      []int  `form:"ids" json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
		Pid      string `form:"pid" json:"pid" default:"" validate:"required,numeric" label:"父级ID"`
		Name     string `form:"name" json:"name" default:"" validate:"omitempty,max=20" label:"名称"`
		Describe string `form:"describe" json:"describe" default:"" validate:"omitempty,max=100" label:"描述"`
		Value    string `form:"value" json:"value" default:"" validate:"omitempty" label:"配置内容"`
		Control  string `form:"control" json:"control" default:"" validate:"omitempty,max=999" label:"控件类型"`
		Sort     string `form:"sort" json:"sort" default:"" validate:"omitempty,numeric" label:"排序"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
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
	if form.Control != "" {
		data["control"] = form.Control
	}
	if form.Sort != "" {
		data["sort"] = form.Sort
	}
	if form.Pid != "" {
		data["pid"] = form.Pid
	}
	err := service2.SystemConfig.UpdateByIds(form.Ids, data)
	if err == nil {
		service2.SystemLog.CreateByQueue(c, system_log.EventEdit, "修改配置", "修改配置信息")
		c.Success("修改数据成功")
	} else {
		c.Error("修改数据失败")
	}
}

// Create 添加信息
func (ControllerConfig) Create(c *sinking_server.Context) {
	type Form struct {
		Pid      int    `form:"pid" json:"pid" default:"" validate:"numeric" label:"父级ID"`
		Name     string `form:"name" json:"name" validate:"required,max=20" label:"名称"`
		Describe string `form:"describe" json:"describe" validate:"required,max=100" label:"描述"`
		Key      string `form:"key" json:"key" default:"" validate:"required,max=50" label:"配置标识"`
		Value    string `form:"value" json:"value" default:"" validate:"required,max=500" label:"配置内容"`
		Control  string `form:"control" json:"control" default:"" validate:"required,max=500" label:"控件类型"`
		Sort     int    `form:"sort" json:"sort" default:"" validate:"omitempty,numeric" label:"排序"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if num, err := service2.SystemConfig.CountByKey(form.Key); err != nil || num > 0 {
		c.Error("该配置项名称已存在")
		return
	}
	err := service2.SystemConfig.Create(&model.SystemConfig{
		Pid:      form.Pid,
		Name:     form.Name,
		Describe: form.Describe,
		Key:      form.Key,
		Value:    form.Value,
		Control:  form.Control,
		Sort:     form.Sort,
	})
	if err == nil {
		go service.GetSystemConfigServiceIns().ClearCache()
		service2.SystemLog.CreateByQueue(c, system_log.EventCreate, "添加配置", "添加配置信息")
		c.Success("添加数据成功")
	} else {
		c.Error("添加数据失败")
	}
}

// Delete 删除信息
func (ControllerConfig) Delete(c *sinking_server.Context) {
	type Form struct {
		Ids []int `form:"ids" json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	err := service2.SystemConfig.DeleteByIds(form.Ids)
	if err == nil {
		go service.GetSystemConfigServiceIns().ClearCache()
		service2.SystemLog.CreateByQueue(c, system_log.EventDelete, "删除配置", "删除配置信息")
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}

// Updates 修改信息
func (ControllerConfig) Updates(c *sinking_server.Context) {
	type Config struct {
		Pid      string `form:"pid" json:"pid" default:"" validate:"omitempty,numeric" label:"父级ID"`
		Name     string `form:"name" json:"name" default:"" validate:"omitempty,max=20" label:"名称"`
		Describe string `form:"describe" json:"describe" default:"" validate:"omitempty,max=100" label:"描述"`
		Key      string `form:"key" json:"key" default:"" validate:"required,max=50" label:"配置标识"`
		Value    string `form:"value" json:"value" default:"" validate:"omitempty,max=999" label:"配置内容"`
		Control  string `form:"control" json:"control" default:"" validate:"omitempty,max=999" label:"控件类型"`
		Sort     string `form:"sort" json:"sort" default:"" validate:"omitempty,numeric" label:"排序"`
	}
	type Form struct {
		Configs []*Config `form:"configs" json:"configs" default:"" validate:"gte=1" label:"配置标识"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	for _, v := range form.Configs {
		data := make(map[string]interface{})
		if v.Name != "" {
			data["name"] = v.Name
		}
		if v.Describe != "" {
			data["describe"] = v.Describe
		}
		if v.Value != "" {
			data["value"] = v.Value
		}
		if v.Control != "" {
			data["control"] = v.Control
		}
		if v.Sort != "" {
			data["sort"] = v.Sort
		}
		if v.Pid != "" {
			data["pid"] = v.Pid
		}
		_ = service2.SystemConfig.UpdateByKey(v.Key, data)
	}
	service2.SystemLog.CreateByQueue(c, system_log.EventEdit, "修改配置", "修改配置信息")
	go service.GetSystemConfigServiceIns().ClearCache()
	c.Success("修改数据成功")
}
