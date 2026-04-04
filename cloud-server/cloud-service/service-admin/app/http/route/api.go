package route

import (
	"service-admin/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApi 加载API
func loadApi(s *sinking_web.Engine) {
	page := s.Group("/page")
	{
		page.ANY("/menu", context.HandleFunc(controller.Page.Menu))       //菜单列表
		page.ANY("/control", context.HandleFunc(controller.Page.Control)) //控件列表
	}
	person := s.Group("/person")
	{
		person.ANY("/info", context.HandleFunc(controller.Person.Info))     //个人信息
		person.ANY("/update", context.HandleFunc(controller.Person.Update)) //修改资料
	}
	security := s.Group("/security")
	{
		security.ANY("/verify", context.HandleFunc(controller.Security.Verify))     //安全验证
		security.ANY("/password", context.HandleFunc(controller.Security.Password)) //修改密码
		security.ANY("/email", context.HandleFunc(controller.Security.Email))       //修改邮箱
		security.ANY("/phone", context.HandleFunc(controller.Security.Phone))       //修改手机
		security.ANY("/account", context.HandleFunc(controller.Security.Account))   //修改账号
	}

	service.Auth.LoadAdminRouteWithRegister(s.Group("/"), func(ANY func(string, string, func(*context.Context))) {
		ANY("/permission/list", "权限列表", controller.Permission.List)
		ANY("/permission/create", "创建权限", controller.Permission.Create)
		ANY("/permission/update", "修改权限", controller.Permission.Update)
		ANY("/permission/delete", "删除权限", controller.Permission.Delete)

		ANY("/menu/list", "菜单列表", controller.Menu.List)
		ANY("/menu/info", "菜单或控件详情", controller.Menu.Info)
		ANY("/menu/controls", "查看页面控件列表", controller.Menu.Controls)
		ANY("/menu/create", "创建列表或控件", controller.Menu.Create)
		ANY("/menu/update", "修改列表或控件", controller.Menu.Update)
		ANY("/menu/delete", "删除列表或控件", controller.Menu.Delete)
	})
}
