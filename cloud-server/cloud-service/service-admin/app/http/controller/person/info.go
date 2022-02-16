package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
)

func (ControllerPerson) Info(c *sinking_server.Context) {
	user := c.GetSystemUserInfo()
	c.SuccessWithData("获取用户信息成功", Info{
		Account:   user.Account,
		NickName:  user.NickName,
		Avatar:    user.Avatar,
		Email:     user.Email,
		Phone:     user.Phone,
		LoginIp:   user.LoginIp,
		LoginTime: user.LoginTime,
	})
}
