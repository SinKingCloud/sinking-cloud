package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/rbac"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Rbac rbac权限控制
func Rbac() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		rpcCheck := c.Request.Header.Get(util.RpcHeaderTokenName) == util2.Conf.GetString(config.ConfRegisterTokenValue)
		if !rpcCheck {
			//获取用户信息
			userInfo := c.GetSystemUserInfo()
			//判断权限信息
			if rbac.GetIns().Check(userInfo.Id, userInfo.RoleId, c.Request.URL.Path, c.Request.Method) {
				c.Next()
			} else {
				c.NotAllow("您的帐户暂无权限操作", nil)
				c.Abort()
				return
			}
		}
	})
}
