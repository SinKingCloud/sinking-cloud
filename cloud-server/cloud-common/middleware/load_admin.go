package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	sinking_web "github.com/SinKingCloud/sinking-go/sinking-web"
)

// LoadAdmin 加载管理员信息
func LoadAdmin() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		c.Next()
	})
}
