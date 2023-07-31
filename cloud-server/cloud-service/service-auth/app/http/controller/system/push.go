package system

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
)

// Push 生成推送token
func (ControllerSystem) Push(c *sinking_server.Context) {
	u := c.GetSystemUserInfo()
	if u == nil || u.Id <= 0 {
		c.Error("获取账户信息失败")
		return
	}
	id, e := sinking.NewWebSocketProducer().GenTokenByUserId(u.Id, true)
	if e != nil {
		c.Error(e.Error())
		return
	}
	c.SuccessWithData("注册推送实例成功", id)
}
