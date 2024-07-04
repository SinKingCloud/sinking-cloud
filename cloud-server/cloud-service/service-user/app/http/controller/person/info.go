package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/service"
	"strconv"
)

func (ControllerPerson) Info(c *sinking_server.Context) {
	user := service.User.GetCurrentUserInfo(c.GetWebUserInfo().Id, false)
	c.SuccessWithData("获取用户信息成功", Info{
		Id:        user.Id,
		Account:   user.Account,
		NickName:  user.NickName,
		Avatar:    user.Avatar,
		Money:     user.Money,
		Email:     user.Email,
		Phone:     str.NewStringTool(strconv.FormatInt(user.Phone, 10)).Show(3, 4, "*"),
		RealAuth:  user.RealAuth,
		LoginIp:   user.LoginIp,
		LoginTime: user.LoginTime,
	})
}
