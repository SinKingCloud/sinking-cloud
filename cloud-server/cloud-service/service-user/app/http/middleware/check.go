package middleware

import (
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// CheckLogin 判断是否登陆
func CheckLogin(c *context.Context) {
	//RPC放行
	if rpc.IsRpcPath(c.Request.URL.Path) {
		c.Next()
		return
	}
	//登陆验证
	token := c.Request.Header.Get(constant.JwtTokenName)
	if token == "" {
		c.NotLogin("您还未登陆,请先登陆账户")
		c.Abort()
		return
	}
	user, err := service.Auth.GetTokenInfo(token)
	if err != nil || user == nil {
		c.TokenError("登陆超时,请重新登陆")
		c.Abort()
		return
	}
	c.SetUserInfo(user)
	c.Next()
}
