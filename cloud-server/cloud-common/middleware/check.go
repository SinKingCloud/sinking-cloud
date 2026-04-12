package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// CheckLogin 判断是否登陆
func CheckLogin(c *context.Context) {
	//RPC请求直接放行
	if rpc.IsRpcPath(c.Request.URL.Path) {
		c.Next()
		return
	}
	u, err := user.GetIns().GetJwtUserInfo(c)
	if err != nil {
		c.NotLogin(err.Error())
		c.Abort()
		return
	}
	c.SetUserInfo(u)
	c.Next()
}
