package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	user2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/service/user"
)

func (ControllerPerson) Info(c *sinking_server.Context) {
	user := user2.GetIns().GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	c.SuccessWithData("获取用户信息成功", Info{
		Id:        user.Id,
		Account:   user.Account,
		NickName:  user.NickName,
		Avatar:    user.Avatar,
		Money:     user.Money,
		Email:     user.Email,
		Phone:     user.Phone,
		RealAuth:  user.RealAuth,
		LoginIp:   user.LoginIp,
		LoginTime: user.LoginTime,
	})
}
