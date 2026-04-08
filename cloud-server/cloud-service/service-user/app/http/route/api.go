package route

import (
	"service-user/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApi 加载API
func loadApi(s *sinking_web.Engine) {
	person := s.Group("/person")
	{
		person.ANY("/info", context.HandleFunc(controller.Person.Info))     //个人信息
		person.ANY("/update", context.HandleFunc(controller.Person.Update)) //修改资料
	}
	security := s.Group("/security")
	{
		security.ANY("/session", context.HandleFunc(controller.Security.Session))   //会话列表
		security.ANY("/logout", context.HandleFunc(controller.Security.Logout))     //注销会话
		security.ANY("/verify", context.HandleFunc(controller.Security.Verify))     //安全验证
		security.ANY("/password", context.HandleFunc(controller.Security.Password)) //修改密码
		security.ANY("/email", context.HandleFunc(controller.Security.Email))       //修改邮箱
		security.ANY("/phone", context.HandleFunc(controller.Security.Phone))       //修改手机
		security.ANY("/account", context.HandleFunc(controller.Security.Account))   //修改账号
	}
}
