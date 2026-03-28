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
		person.ANY("/info", context.HandleFunc(controller.Person.Info))         //个人信息
		person.ANY("/password", context.HandleFunc(controller.Person.Password)) //修改密码
		person.ANY("/email", context.HandleFunc(controller.Person.Email))       //修改邮箱
		person.ANY("/phone", context.HandleFunc(controller.Person.Phone))       //修改手机
		person.ANY("/update", context.HandleFunc(controller.Person.Update))     //修改资料
	}
}
