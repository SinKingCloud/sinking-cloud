package route

import (
	"service-user/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadUser 加载user
func loadUser(s *sinking_web.Engine) {
	g := service.Auth.GetUserRouteGroup(s)
	person := g.Group("/person")
	{
		person.ANY("/info", context.HandleFunc(controller.User.Person.Info))     //个人信息
		person.ANY("/update", context.HandleFunc(controller.User.Person.Update)) //修改资料
	}
	security := g.Group("/security")
	{
		security.ANY("/certificate", context.HandleFunc(controller.User.Security.Certificate)) //实名认证
		security.ANY("/session", context.HandleFunc(controller.User.Security.Session))         //会话列表
		security.ANY("/logout", context.HandleFunc(controller.User.Security.Logout))           //注销会话
		security.ANY("/verify", context.HandleFunc(controller.User.Security.Verify))           //安全验证
		security.ANY("/password", context.HandleFunc(controller.User.Security.Password))       //修改密码
		security.ANY("/email", context.HandleFunc(controller.User.Security.Email))             //修改邮箱
		security.ANY("/phone", context.HandleFunc(controller.User.Security.Phone))             //修改手机
		security.ANY("/account", context.HandleFunc(controller.User.Security.Account))         //修改账号
	}
}
