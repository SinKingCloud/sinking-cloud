package message

import (
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Clear 清空消息
func (Controller) Clear(c *context.Context) {
	if err := service.Message.Clear(c.GetUserInfo().Id); err != nil {
		c.Error(err.Error())
		return
	}
	c.Success("清空消息成功")
}
