package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// RpcCheck 微服务间调用安全判断
func RpcCheck() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		if c.Request.Header.Get(util2.RpcHeaderTokenName) != util.Conf.GetString(config.ConfRegisterTokenValue) {
			c.ServiceError("RPC安全认证失败")
			c.Abort()
		} else {
			c.Next()
		}
	})
}
