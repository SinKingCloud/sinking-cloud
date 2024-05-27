package system

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service"
)

// OutLogin 退出登录
func (ControllerSystem) OutLogin(c *sinking_server.Context) {
	user, err := service2.GetSystemUserServiceIns().FindByIdWithCache(c.GetSystemUserInfo().Id, true)
	if err == nil && user != nil {
		service.Auth.OutLogin(c.Request.Header.Get(util.JwtTokenNameType), user, nil)
	}
	rpc.Log.CreateAdmin(user.Id, c.GetRequestId(), c.GetRequestIp(), log.UserEventLogin, "注销登录", "注销登陆成功")
	c.Success("注销登录成功")
}
