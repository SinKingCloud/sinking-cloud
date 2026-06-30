package message

import (
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Clear 清空消息
func (Controller) Clear(c *context.Context) {
	u := c.GetUserInfo()
	if err := service.Message.Clear(u.Id); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("清空消息成功")
}
