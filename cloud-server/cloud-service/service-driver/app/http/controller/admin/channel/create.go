package channel

import (
	"service-driver/app/service"

	commonService "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user/user_log_type"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Create 创建数据
func (Controller) Create(c *context.Context) {
	var form struct {
		Name   string `json:"name" default:"" validate:"required" label:"通道名称"`
		Type   int    `json:"type" default:"0" validate:"numeric" label:"通道类型"`
		Status int    `json:"status" default:"0" validate:"numeric" label:"通道状态"`
		Config string `json:"config" default:"" validate:"omitempty" label:"通道配置"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	err := service.Channel.Create(form.Name, form.Type, form.Status, form.Config)
	if err != nil {
		c.Error(err.Error())
		return
	}
	_ = commonService.User.CreateSystemUserLog(c.GetUserInfo().Id, c.GetRequestId(), c.GetRequestIp(), user_log_type.EventCreate, "创建存储通道", "创建存储通道["+form.Name+"]")
	c.Success("创建数据成功")
}
