package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"time"
)

// AuthAdmin 后台系统鉴权
func AuthAdmin() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		rpcCheck := c.Request.Header.Get(util.RpcHeaderTokenName) == util2.Conf.GetString(config.ConfRegisterTokenValue)
		if !rpcCheck {
			token := c.Request.Header.Get(util.JwtTokenName)
			if token == "" {
				c.NotLogin("您还未登陆,请登陆账户", nil)
				c.Abort()
				return
			}
			key := jwt.CheckToken(token)
			if key == nil || time.Now().Unix() > key.ExpiresAt || key.SystemUser == nil {
				c.TokenError("登陆超时,请重新登陆", nil)
				c.Abort()
				return
			}
			c.SetSystemUserInfo(key.SystemUser)
		}
		c.Next()
	})
}
