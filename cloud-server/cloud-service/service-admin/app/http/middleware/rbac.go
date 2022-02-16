package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/rbac"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Rbac rbac权限控制
func Rbac() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		//获取用户信息
		userInfo := c.GetSystemUserInfo()
		//判断权限信息
		if rbac.CheckAccessByCache(userInfo.RoleId, c.Request.RequestURI, c.Request.Method) {
			c.Next()
		} else {
			c.NotAllow("您的帐户暂无权限操作", nil)
			c.Abort()
		}
	})
}
