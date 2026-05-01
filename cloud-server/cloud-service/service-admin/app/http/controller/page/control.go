package page

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Control 页面控件
func (Controller) Control(c *context.Context) {
	var form struct {
		Code string `json:"code" default:"" validate:"required" label:"页面标识"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	controls := service.Rbac.GetPageControlsByUserId(c.GetUserInfo().Id, form.Code)
	if controls == nil {
		c.Error("获取控件数据失败")
		return
	}
	c.SuccessWithData("获取控件数据成功", controls)
}
