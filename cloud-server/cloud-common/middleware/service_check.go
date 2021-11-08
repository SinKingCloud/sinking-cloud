package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ServiceCheck 微服务间通信安全判断
func ServiceCheck() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		if c.Request.Header.Get(util.Conf.GetString("register.tokenKey")) != util.Conf.GetString("register.tokenValue") {
			c.ServiceError("微服务安全认证失败")
			c.Abort()
		} else {
			c.Next()
		}
	})
}
