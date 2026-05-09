package middleware

import (
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
)

// Rbac rbac权限拦截
func Rbac(c *context.Context) {
	user := c.GetUserInfo()
	if user != nil && user.Id > 0 {
		err := service.Rbac.CheckUserHasPermission(user.Id, "/"+constant.ServiceAdmin+c.Request.URL.Path)
		if err != nil {
			c.NotAllow("暂无权限访问此内容")
			c.Abort()
			return
		}
	}
	c.Next()
}
