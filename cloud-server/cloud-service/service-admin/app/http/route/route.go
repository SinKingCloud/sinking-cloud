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
	role.GET("/list", sinking_server.HandleFunc(roleHandle.List))                       //角色列表
	role.POST("/update", sinking_server.HandleFunc(roleHandle.Update))                  //角色更新
	role.POST("/create", sinking_server.HandleFunc(roleHandle.Create))                  //角色添加
	role.POST("/delete", sinking_server.HandleFunc(roleHandle.Delete))                  //角色删除
	role.GET("/reset_access", sinking_server.HandleFunc(roleHandle.ResetAllRoleAccess)) //更新角色权限
	//用户管理
	userHandle := &systems.ControllerUser{}
	user := r.Group("/user")
	user.GET("/list", sinking_server.HandleFunc(userHandle.List))      //用户列表
	user.POST("/update", sinking_server.HandleFunc(userHandle.Update)) //用户更新
	user.POST("/create", sinking_server.HandleFunc(userHandle.Create)) //用户添加
	user.POST("/delete", sinking_server.HandleFunc(userHandle.Delete)) //用户删除
	//配置管理
	configHandle := &systems.ControllerConfig{}
	config := r.Group("/config")
	config.GET("/list", sinking_server.HandleFunc(configHandle.List))        //配置列表
	config.POST("/update", sinking_server.HandleFunc(configHandle.Update))   //配置更新
	config.POST("/updates", sinking_server.HandleFunc(configHandle.Updates)) //配置更新(批量)
	config.POST("/create", sinking_server.HandleFunc(configHandle.Create))   //配置添加
	config.POST("/delete", sinking_server.HandleFunc(configHandle.Delete))   //配置删除
	//日志管理
	logHandle := &systems.ControllerLog{}
	log := r.Group("/log")
	log.GET("/list", sinking_server.HandleFunc(logHandle.List))      //配置列表
	log.POST("/delete", sinking_server.HandleFunc(logHandle.Delete)) //配置更新
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
	r.POST("/update", sinking_server.HandleFunc(handle.Update))
	r.POST("/password", sinking_server.HandleFunc(handle.Password))
}
