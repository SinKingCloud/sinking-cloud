package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

func Login(c *context.Context) {}

// loadApi 加载API
func loadAdmin(s *sinking_web.Engine) {
	service.Auth.LoadAdminRoute(s, func(ANY func(string, string, func(*context.Context))) {
		ANY("/test", "测试权限", Login)
		ANY("/test2", "测试权限2", Login)
	})
}
