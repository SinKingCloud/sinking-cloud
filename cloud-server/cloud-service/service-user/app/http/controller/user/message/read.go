package message

import (
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Read 消息已读
func (Controller) Read(c *context.Context) {
	var form struct {
		Ids []int64 `json:"ids" default:"" validate:"omitempty,min=1,max=1000,unique" label:"消息ID列表"`
		All bool    `json:"all" default:"false" label:"是否全部已读"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	if (form.Ids == nil || len(form.Ids) <= 0) && !form.All {
		c.Error("参数不能为空")
	}
	u := c.GetUserInfo()
	if form.All {
		if err := service.Message.ReadAll(u.Id); err != nil {
			c.Error(err.Error())
			return
		}
		c.Success("操作成功")
		return
	}
	if err := service.Message.Read(u.Id, form.Ids); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("操作成功")
}
