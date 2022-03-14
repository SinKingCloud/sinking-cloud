package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/http/controller/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/http/controller/person"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/http/controller/systems"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	loadPageRoute(s)
	loadPersonRoute(s)
	loadSystemsRoute(s)
}

// loadSystemRoute 加载系统管理路由
func loadSystemsRoute(s *sinking_web.Engine) {
	r := s.Group("/systems")
	//权限管理
	accessHandle := &systems.ControllerAccess{}
	access := r.Group("/access")
	access.GET("/list", sinking_server.HandleFunc(accessHandle.List))      //权限列表
	access.POST("/update", sinking_server.HandleFunc(accessHandle.Update)) //权限修改
	access.POST("/create", sinking_server.HandleFunc(accessHandle.Create)) //权限添加
	access.POST("/delete", sinking_server.HandleFunc(accessHandle.Delete)) //权限删除
	//菜单管理
	menuHandle := &systems.ControllerMenu{}
	menu := r.Group("/menu")
	menu.GET("/list", sinking_server.HandleFunc(menuHandle.List))      //菜单列表
	menu.POST("/update", sinking_server.HandleFunc(menuHandle.Update)) //菜单更新
	menu.POST("/create", sinking_server.HandleFunc(menuHandle.Create)) //菜单添加
	menu.POST("/delete", sinking_server.HandleFunc(menuHandle.Delete)) //菜单删除
	//请求日志
	logRequestHandle := &systems.ControllerLogRequest{}
	logRequest := r.Group("/log_request")
	logRequest.GET("/list", sinking_server.HandleFunc(logRequestHandle.List))
	logRequest.GET("/info", sinking_server.HandleFunc(logRequestHandle.Info))
	logRequest.POST("/delete", sinking_server.HandleFunc(logRequestHandle.Delete))
	//角色管理
	roleHandle := &systems.ControllerRole{}
	role := r.Group("/role")
	role.GET("/list", sinking_server.HandleFunc(roleHandle.List))      //角色列表
	role.POST("/update", sinking_server.HandleFunc(roleHandle.Update)) //角色更新
	role.POST("/create", sinking_server.HandleFunc(roleHandle.Create)) //角色添加
	role.POST("/delete", sinking_server.HandleFunc(roleHandle.Delete)) //角色删除
}

// loadPageRoute 加载页面操作路由
func loadPageRoute(s *sinking_web.Engine) {
	r := s.Group("/page")
	handle := &page.ControllerPage{}
	r.GET("/menu", sinking_server.HandleFunc(handle.Menu))       //获取菜单列表
	r.GET("/control", sinking_server.HandleFunc(handle.Control)) //获取页面控件
}

// loadPersonRoute 加载个人操作路由
func loadPersonRoute(s *sinking_web.Engine) {
	r := s.Group("/person")
	handle := &person.ControllerPerson{}
	r.GET("/info", sinking_server.HandleFunc(handle.Info))
}
