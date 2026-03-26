package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// CheckLogin 判断是否登陆
func CheckLogin(c *context.Context) {
	//RPC请求直接放行
	if rpc.IsRpcPath(c.Request.URL.Path) {
		c.Next()
		return
	}
	token := c.Request.Header.Get(constant.JwtTokenName)
	if token == "" {
		c.NotLogin("您还未登陆,请先登陆账户")
		c.Abort()
		return
	}
	key := jwt.CheckToken(token)
	if key == nil || key.User == nil || key.User.Id <= 0 || key.User.Token == "" {
		c.TokenError("登陆超时,请重新登陆")
		c.Abort()
		return
	}
	c.SetUserInfo(key.User)
	c.Next()
}
