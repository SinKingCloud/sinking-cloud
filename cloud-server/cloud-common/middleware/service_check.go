package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ServiceCheck 微服务间通信安全判断
func ServiceCheck() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		c.Next()
	})
}
