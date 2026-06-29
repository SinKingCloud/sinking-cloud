package message

import (
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Delete 删除消息
func (Controller) Delete(c *context.Context) {
	var form struct {
		Ids []int64 `json:"ids" default:"" validate:"required,min=1,max=1000,unique" label:"消息ID列表"`
	}
	if ok, msg := c.ValidatorAll(&form); !ok {
		c.Error(msg)
		return
	}
	u := c.GetUserInfo()
	if err := service.Message.Delete(u.Id, form.Ids, u.Token); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("删除数据成功")
}
