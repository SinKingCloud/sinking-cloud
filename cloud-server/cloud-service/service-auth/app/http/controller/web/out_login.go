package web

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-auth/app/service/auth"
)

// OutLogin 退出登录
func (ControllerWeb) OutLogin(c *sinking_server.Context) {
	user, err := service2.GetWebUserServiceIns().FindByIdWithCache(c.GetWebUserInfo().Id, true)
	if err == nil && user != nil {
		auth.OutLogin(c.Request.Header.Get(util.JwtTokenNameType), nil, user)
	}
	c.Success("注销登录成功")
}
