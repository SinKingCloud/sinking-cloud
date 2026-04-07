package menu

import (
	"errors"
	menu2 "service-admin/app/repository/menu"
	"service-admin/app/service"
	"strconv"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

type Controller struct {
}

// List 菜单列表
func (Controller) List(c *context.Context) {
	list := service.Rbac.GetAllMenuTreeByDirAndPage()
	c.SuccessWithData("获取数据成功", list)
}

// Info 页面信息
func (Controller) Info(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"记录ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	info, err := service.Rbac.GetMenuOrControlInfoById(form.Id)
	if err != nil {
		c.Error(errors.Join(errors.New("获取数据失败"), err))
		return
	}
	c.SuccessWithData("获取数据成功", info)
}

// Controls 控件列表
func (Controller) Controls(c *context.Context) {
	var form struct {
		Id int64 `json:"id" default:"" validate:"numeric,min=1" label:"页面ID"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	list := service.Rbac.GetPageControlsByMenuId(form.Id)
	c.SuccessWithData("获取数据成功", list)
}

// Create 创建菜单或控件
func (Controller) Create(c *context.Context) {
	var form struct {
		ParentId        int64    `json:"parent_id" default:"0" validate:"numeric,gte=0" label:"父级"`
		Code            string   `json:"code" default:"" validate:"required,max=100" label:"标识"`
		Name            string   `json:"name" default:"" validate:"required,max=100" label:"名称"`
		Icon            string   `json:"icon" default:"" validate:"omitempty,max=30" label:"图标"`
		Type            int      `json:"type" default:"0" validate:"numeric" label:"类型"`
		Sort            int      `json:"sort" default:"0" validate:"numeric" label:"排序"`
		Status          int      `json:"status" default:"0" validate:"numeric" label:"状态"`
		PermissionCodes []string `json:"permission_codes" default:"" validate:"omitempty,unique" label:"权限标识"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Rbac.CreateMenuOrControl(form.ParentId, form.Code, form.Name, form.Icon, form.Type, form.Sort, form.Status, form.PermissionCodes)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("创建数据成功")
}

// Update 更新菜单或控件
func (Controller) Update(c *context.Context) {
	var form struct {
		Id              int64     `json:"id" validate:"required,min=1" label:"记录ID"`
		ParentId        string    `json:"parent_id" validate:"omitempty,numeric,gte=0" label:"父级"`
		Name            string    `json:"name" validate:"omitempty,max=100" label:"名称"`
		Icon            string    `json:"icon" validate:"omitempty,max=30" label:"图标"`
		Sort            string    `json:"sort" validate:"omitempty,numeric" label:"排序"`
		Status          string    `json:"status" validate:"omitempty,numeric" label:"状态"`
		PermissionCodes *[]string `json:"permission_codes" validate:"omitempty,unique" label:"权限标识"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	data := &menu2.UpdateMenu{}
	if form.ParentId != "" {
		data.ParentId, _ = strconv.ParseInt(form.ParentId, 10, 64)
	}
	if form.Name != "" {
		data.Name = form.Name
	}
	if form.Icon != "" {
		data.Icon = form.Icon
	}
	if form.Sort != "" {
		data.Sort, _ = strconv.Atoi(form.Sort)
	}
	if form.Status != "" {
		data.Status, _ = strconv.Atoi(form.Status)
	}
	err := service.Rbac.UpdateMenuOrControl(form.Id, data, form.PermissionCodes)
	if err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("更新数据成功")
}

// Delete 删除菜单或控件
func (Controller) Delete(c *context.Context) {
	var form struct {
		Ids []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"ID列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Rbac.DeleteMenuOrControlByIds(form.Ids)
	if err != nil {
		c.Error(errors.Join(errors.New("删除数据失败"), err))
	} else {
		c.Success("删除数据成功")
	}
}
