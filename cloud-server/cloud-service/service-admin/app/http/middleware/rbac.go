package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/log"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/service/rbac"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Rbac rbac权限控制
func Rbac() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		//获取用户信息
		userInfo := c.GetSystemUserInfo()
		//判断权限信息
		if userInfo.Id == constant.AdminId {
			c.Next()
			return
		}
		if r, d := rbac.GetIns().CheckAccessByCache(userInfo.RoleId, c.Request.URL.Path, c.Request.Method); r {
			if d != nil {
				types := 0
				switch d.Method {
				case "GET":
					types = 0
					break
				case "POST":
					types = 1
					break
				case "UPDATE":
					types = 2
					break
				case "DELETE":
					types = 3
					break
				}
				log.GetIns().CreateByQueue(
					c.GetSystemUserInfo().Id,
					c.GetRequestId(),
					c.GetRequestIp(),
					types,
					d.Name,
					d.Describe,
				)
			}
			c.Next()
		} else {
			c.NotAllow("您的帐户暂无权限操作", nil)
			c.Abort()
		}
	})
}
