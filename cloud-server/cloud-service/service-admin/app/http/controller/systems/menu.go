package systems

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_menu"
)

// List 获取菜单数据
func (ControllerMenu) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField string `form:"order_by_field" json:"order_by_field" default:"sort" validate:"oneof=id sort update_time create_time status" label:"排序字段"`
		OrderByType  string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Title        string `form:"title" default:"" validate:"omitempty,max=100" label:"标题"`
		Key          string `form:"key" default:"" validate:"omitempty,max=50" label:"标识"`
		Path         string `form:"path" default:"" validate:"omitempty,max=50" label:"目录"`
		Type         string `form:"type" default:"" validate:"omitempty,oneof=0 1" label:"类型"`
		Status       string `form:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorGet(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.Title != "" {
		where["title"] = form.Title
	}
	if form.Key != "" {
		where["key"] = form.Key
	}
	if form.Path != "" {
		where["path"] = form.Path
	}
	if form.Status != "" {
		where["status"] = form.Status
	}
	if form.Type != "" {
		where["type"] = form.Type
	}
	data, total, err := system_menu.GetIns().Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventShow, "查看菜单", "查看菜单列表")
		c.Error("获取菜单数据失败")
	} else {
		c.SuccessWithData("获取菜单数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}

// Update 修改菜单数据
func (ControllerMenu) Update(c *sinking_server.Context) {
	type Form struct {
		Ids      []int  `json:"ids" default:"" validate:"required,min=1,unique" label:"ID列表"`
		Title    string `json:"title" default:"" validate:"omitempty,max=100" label:"标题"`
		Key      string `json:"key" default:"" validate:"omitempty,max=50" label:"标识"`
		Path     string `json:"path" default:"" validate:"omitempty,max=50" label:"目录"`
		Icon     string `json:"icon" default:"" validate:"omitempty" label:"图标"`
		Accesses string `json:"accesses" default:"" validate:"omitempty" label:"关联权限"`
		Sort     string `json:"sort" default:"" validate:"omitempty,numeric" label:"排序"`
		Status   string `json:"status" default:"" validate:"omitempty,numeric,oneof=0 1" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	data := make(map[string]interface{})
	if form.Title != "" {
		data["title"] = form.Title
	}
	if form.Key != "" {
		data["key"] = form.Key
	}
	if form.Path != "" {
		data["path"] = form.Path
	}
	if form.Icon != "" {
		data["icon"] = form.Icon
	}
	if form.Accesses != "" {
		data["accesses"] = form.Accesses
	}
	if form.Sort != "" {
		data["sort"] = form.Sort
	}
	if form.Status != "" {
		data["status"] = form.Status
	}
	err := system_menu.GetIns().UpdateByIds(form.Ids, data)
	if err == nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventEdit, "修改菜单", "修改菜单信息")
		c.Success("修改数据成功")
	} else {
		c.Error("修改数据失败")
	}
}

// Create 添加菜单数据
func (ControllerMenu) Create(c *sinking_server.Context) {
	type Form struct {
		Pid      int    `json:"pid" default:"" validate:"numeric" label:"父级id"`
		Type     int    `json:"type" default:"" validate:"numeric,oneof=0 1" label:"类型"`
		Title    string `json:"title" default:"" validate:"required,max=100" label:"标题"`
		Key      string `json:"key" default:"" validate:"required,max=50" label:"标识"`
		Path     string `json:"path" default:"" validate:"omitempty,max=50" label:"目录"`
		Icon     string `json:"icon" default:"" validate:"omitempty" label:"图标"`
		Accesses string `json:"accesses" default:"" validate:"omitempty" label:"关联权限"`
		Sort     int    `json:"sort" default:"" validate:"omitempty,numeric" label:"排序"`
		Status   int    `json:"status" default:"" validate:"numeric,oneof=0 1" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := system_menu.GetIns().Create(&model.SystemMenu{
		Pid:      form.Pid,
		Title:    form.Title,
		Type:     form.Type,
		Key:      form.Key,
		Path:     form.Path,
		Icon:     form.Icon,
		Accesses: form.Accesses,
		Status:   form.Status,
		Sort:     form.Sort,
	})
	if err == nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventCreate, "添加菜单", "添加菜单信息")
		c.Success("添加数据成功")
	} else {
		c.Error("添加数据失败")
	}
}

// Delete 添加菜单信息
func (ControllerMenu) Delete(c *sinking_server.Context) {
	type Form struct {
		Ids []int `json:"ids" default:"" validate:"required,min=1,unique" label:"ID列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorJson(form); !ok {
		c.Error(msg)
		return
	}
	err := system_menu.GetIns().DeleteByIds(form.Ids)
	if err == nil {
		system_log.GetIns().CreateByQueue(c, system_log.EventDelete, "删除菜单", "删除菜单记录")
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}
