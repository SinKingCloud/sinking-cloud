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

// AuthUser 用户系统鉴权
func AuthUser() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		rpcCheck := c.Request.Header.Get(util.RpcHeaderTokenName) == util2.Conf.GetString(config.ConfRegisterTokenValue)
		if !rpcCheck {
			token := c.Request.Header.Get(util.JwtTokenName)
			if token == "" {
				c.NotLogin("您还未登陆,请先登陆账户", nil)
				c.Abort()
				return
			}
			key := jwt.CheckToken(token)
			if key == nil || time.Now().Unix() > key.ExpiresAt || key.WebUser == nil {
				c.TokenError("您还未登陆,请登陆账户", nil)
				c.Abort()
				return
			}
			c.SetWebUserInfo(key.WebUser)
		}
		c.Next()
	})
}
