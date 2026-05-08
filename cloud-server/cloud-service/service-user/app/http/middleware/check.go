package middleware

import (
	"service-user/app/enum/user_status"
	"service-user/app/service"

	service2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// CheckLogin 判断是否登陆
func CheckLogin(c *context.Context) {
	//RPC放行
	if rpc.IsRpcPath(c.Request.URL.Path) {
		c.Next()
		return
	}
	var u *jwt.User
	var err error
	u = c.GetUserInfo()
	if u == nil || u.Id <= 0 {
		u, err = service2.User.GetJwtUserInfo(c)
		if err != nil {
			c.NotLogin("登陆超时,请重新登陆")
			c.Abort()
			return
		}
	}
	session := service.Auth.GetSessionInfo(u.Token)
	if session == nil {
		c.NotLogin("登陆超时,请重新登陆")
		c.Abort()
		return
	}
	user, err := service.User.GetInfo(u.Id, false)
	if err != nil || user == nil || user.Id <= 0 {
		c.NotLogin("登陆超时,请重新登陆")
		c.Abort()
		return
	}
	if user.Status != user_status.Enable {
		c.NotLogin("该账户已被封禁,如有问题请咨询官方客服")
		c.Abort()
		return
	}
	c.SetUserInfo(u)
	c.Next()
}
