package person

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	role2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/system_role"
	user2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/system_user"
)

func (ControllerPerson) Info(c *sinking_server.Context) {
	user := user2.GetIns().GetCurrentUserInfo(c.GetSystemUserInfo().Id, false)
	role, _ := role2.GetIns().FindByIdWithCache(user.RoleId, false)
	c.SuccessWithData("获取用户信息成功", Info{
		Id:        user.Id,
		Account:   user.Account,
		NickName:  user.NickName,
		RoleName:  role.Name,
		Avatar:    user.Avatar,
		Email:     user.Email,
		Phone:     user.Phone,
		LoginIp:   user.LoginIp,
		LoginTime: user.LoginTime,
	})
}
