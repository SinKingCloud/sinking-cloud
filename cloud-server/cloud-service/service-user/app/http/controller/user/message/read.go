package message

import (
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Read 消息已读
func (Controller) Read(c *context.Context) {
	var form struct {
		Ids []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"消息ID列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if err := service.Message.Read(c.GetUserInfo().Id, form.Ids); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("操作成功")
}
