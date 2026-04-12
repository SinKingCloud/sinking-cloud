package controller

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Socket 获取socket token
func Socket(c *context.Context) {
	u, err := service.User.GetJwtUserInfo(c)
	if err != nil {
		c.Error(err)
		return
	}
	if u.Id <= 0 || u.Token == "" {
		c.Error("获取会话信息失败")
		return
	}
	token, err := service.Auth.GenWebUserSocketToken(u.Token)
	if err != nil {
		c.Error("获取数据失败")
		return
	}
	c.SuccessWithData("获取数据成功", token)
}
