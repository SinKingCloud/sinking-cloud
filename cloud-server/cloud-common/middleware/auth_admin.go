package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"time"
)

// AuthAdmin 后台系统鉴权
func AuthAdmin() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		token := c.Request.Header.Get(util.JwtTokenNameSystemUser)
		if token == "" {
			c.NotLogin("您还未登陆,请登陆账户", nil)
			c.Abort()
			return
		}
		key := jwt.CheckToken(token)
		if key == nil || time.Now().Unix() > key.ExpiresAt {
			c.TokenError("登陆超时,请重新登陆", nil)
			c.Abort()
			return
		}
		c.Set(util.SystemUserInfo, key.User)
		//判断权限
		c.Next()
	})
}
