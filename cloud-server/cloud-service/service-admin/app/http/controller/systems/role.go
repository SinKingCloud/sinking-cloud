package systems

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/service/system_log"
)

// List 列表
func (ControllerRole) List(c *sinking_server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `form:"order_by_field" json:"order_by_field" default:"id" validate:"oneof=id update_time create_time status" label:"排序字段"`
		OrderByType     string `form:"order_by_type" json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Name            string `form:"name" json:"name" default:"" validate:"max=20" label:"名称"`
		Describe        string `form:"describe" json:"describe" default:"" validate:"max=100" label:"描述"`
		Status          string `form:"status" json:"status" default:"" validate:"omitempty,oneof=0 1" label:"状态"`
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
	if form.Name != "" {
		where["name"] = form.Name
	}
	if form.Describe != "" {
		where["describe"] = form.Describe
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
	data, total, err := service.SystemRole.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err == nil {
		service.SystemLog.CreateByQueue(c, system_log.EventShow, "查看角色", "查看角色列表")
		c.SuccessWithData("获取角色数据成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
		return
	}
	c.Error("获取数据失败")
}

// Update 修改信息
func (r ControllerRole) Update(c *sinking_server.Context) {
	type Form struct {
		Ids      []int  `form:"ids" json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
		Name     string `form:"name" json:"name" default:"" validate:"omitempty,max=20" label:"名称"`
		Describe string `form:"describe" json:"describe" default:"" validate:"omitempty,max=100" label:"描述"`
		Menus    string `form:"menus" json:"menus" default:"" validate:"omitempty" label:"绑定菜单"`
		Status   string `form:"status" json:"status" default:"" validate:"omitempty,numeric,oneof=0 1" label:"状态"`
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
	if form.Menus != "" {
		data["menus"] = form.Menus
		data["accesses"] = r.getRoleIdsBYMenuIds(form.Menus)
		defer func() {
			go r.resetAllRoleAccess()
		}()
	}
	if form.Status != "" {
		data["status"] = form.Status
	}
	err := service.SystemRole.UpdateByIds(form.Ids, data)
	if err == nil {
		service.SystemLog.CreateByQueue(c, system_log.EventEdit, "修改角色", "修改角色信息")
		c.Success("修改数据成功")
	} else {
		c.Error("修改数据失败")
	}
}

// Create 添加信息
func (r ControllerRole) Create(c *sinking_server.Context) {
	type Form struct {
		Name     string `form:"name" json:"name" validate:"required,max=20" label:"名称"`
		Describe string `form:"describe" json:"describe" validate:"required,max=100" label:"描述"`
		Menus    string `form:"menus" json:"menus" default:"" validate:"required" label:"绑定菜单"`
		Status   int    `form:"status" json:"status" validate:"oneof=0 1" label:"状态"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	err := service.SystemRole.Create(&model.SystemRole{
		Name:     form.Name,
		Describe: form.Describe,
		Menus:    form.Menus,
		Accesses: r.getRoleIdsBYMenuIds(form.Menus),
		Status:   form.Status,
	})
	if err == nil {
		service.SystemLog.CreateByQueue(c, system_log.EventCreate, "添加角色", "添加角色信息")
		c.Success("添加数据成功")
	} else {
		c.Error("添加数据失败")
	}
}

// Delete 删除信息
func (r ControllerRole) Delete(c *sinking_server.Context) {
	type Form struct {
		Ids []int `form:"ids" json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if total, err := service.SystemUser.CountByIds(form.Ids); err != nil || total > 0 {
		c.Error("请先删除该角色下的账户")
		return
	}
	err := service.SystemRole.DeleteByIds(form.Ids)
	if err == nil {
		go r.resetAllRoleAccess()
		service.SystemLog.CreateByQueue(c, system_log.EventDelete, "删除角色", "删除角色记录")
		c.Success("删除数据成功")
	} else {
		c.Error("删除数据失败")
	}
}

// getRoleIdsBYMenuIds 获取关联权限ID
func (ControllerRole) getRoleIdsBYMenuIds(menus string) string {
	accessIds := "[]"
	var ids []int
	err := json.Unmarshal([]byte(menus), &ids)
	if err == nil {
		accesses, tempErr := service.SystemMenu.GetAccessIdsByMenuIds(ids)
		if tempErr == nil {
			text, tempErr := json.Marshal(accesses)
			if tempErr == nil {
				accessIds = string(text)
			}
		}
	}
	return accessIds
}

// resetAllRoleAccess 刷新角色关联权限
func (ControllerRole) resetAllRoleAccess() {
	roles, err := service.SystemRole.SelectAll()
	if err == nil {
		for _, v := range roles {
			go func(v *model.SystemRole) {
				_ = service.SystemRole.UpdateById(v.Id, map[string]interface{}{"accesses": v.Menus})
				service.SystemMenu.ClearRoleMenuCache(v.Id)
				service.SystemRole.ClearRoleCache(v.Id)
			}(v)
		}
	}
}

// ResetAllRoleAccess 刷新角色关联权限
func (r ControllerRole) ResetAllRoleAccess(c *sinking_server.Context) {
	go r.resetAllRoleAccess()
	service.SystemLog.CreateByQueue(c, system_log.EventEdit, "刷新角色", "刷新角色权限")
	c.Success("更新权限数据成功")
}
