package middleware

import (
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// AuthAdmin 后台鉴权
func AuthAdmin() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		c.Next()
	})
}

// AuthUser 前台鉴权
func AuthUser() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		fmt.Println("执行成功")
		c.Next()
	})
}
