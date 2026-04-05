package menu

import (
	"errors"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

type Controller struct {
}

// List 菜单列表
func (Controller) List(c *context.Context) {
	list := service.Rbac.GetAllMenuTree()
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

}

// Update 更新菜单或控件
func (Controller) Update(c *context.Context) {

}

// Delete 删除菜单或控件
func (Controller) Delete(c *context.Context) {

}
