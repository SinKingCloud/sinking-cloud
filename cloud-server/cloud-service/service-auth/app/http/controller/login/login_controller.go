package login

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"

func (ControllerLogin) Index(c *sinking_server.Context) {
	c.SuccessWithData("请求成功", c.GetSystemConfig("admin.token.expire"))
}
