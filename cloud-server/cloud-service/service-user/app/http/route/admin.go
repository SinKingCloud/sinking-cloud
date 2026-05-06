package route

import (
	"service-user/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadUser 加载API
func loadAdmin(s *sinking_web.Engine) {
	g := service.Auth.GetAdminRouteGroup(s)
	service.Auth.LoadAdminRouteWithRegister(g, func(ANY func(string, string, func(*context.Context))) {
		ANY("/list", "用户列表数据", controller.Admin.User.List)
		ANY("/update", "修改用户信息", controller.Admin.User.Update)
		ANY("/session", "用户会话列表", controller.Admin.User.Session)
		ANY("/logout", "注销用户会话", controller.Admin.User.Logout)
	})
}
