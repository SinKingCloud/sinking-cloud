package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// LoadConfig 加载系统核心配置
func LoadConfig() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		//加载系统配置
		config := service.GetConfigServiceIns().GetConfigByCache(false)
		c.SetSystemConfig(config)
		c.Next()
	})
}
