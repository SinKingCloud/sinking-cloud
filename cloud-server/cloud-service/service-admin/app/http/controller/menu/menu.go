package menu

import (
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

}

// Controls 控件列表
func (Controller) Controls(c *context.Context) {

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
