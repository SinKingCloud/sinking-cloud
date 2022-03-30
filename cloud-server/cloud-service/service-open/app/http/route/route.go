package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/http/controller/admin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/http/controller/api"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/http/controller/user/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/http/controller/user/uin"
	middleware2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/http/middleware"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	loadApiRoute(s)
	loadUserRoute(s)
}

//加载api路由
func loadApiRoute(s *sinking_web.Engine) {
	g := s.Group("/api")
	g.Use(middleware2.ApiCheck())
	g.GET("/test", sinking_server.HandleFunc(api.Index))
}

// 用户系统路由
func loadUserRoute(s *sinking_web.Engine) {
	g := s.Group("/user")
	g.Use(middleware.AuthUser())
	//用户应用
	handle := &app.ControllerApp{}
	appGroup := g.Group("/app")
	appGroup.GET("/list", sinking_server.HandleFunc(handle.List))
	appGroup.POST("/create", sinking_server.HandleFunc(handle.Create))
	appGroup.POST("/update", sinking_server.HandleFunc(handle.Update))
	appGroup.POST("/delete", sinking_server.HandleFunc(handle.Delete))

	//挂机应用
	uinGroup := g.Group("/uin")
	handleUin := &uin.ControllerUin{}
	uinGroup.GET("/list", sinking_server.HandleFunc(handleUin.List))      //挂机列表
	uinGroup.POST("/delete", sinking_server.HandleFunc(handleUin.Delete)) //删除挂机
	uinGroup.GET("/self", sinking_server.HandleFunc(handleUin.Self))      //概览信息
	uinGroup.GET("/count", sinking_server.HandleFunc(handleUin.Count))    //统计信息
	uinGroup.POST("/buy", sinking_server.HandleFunc(handleUin.Buy))       //业务下单

	//QQ操作
	uinSetGroup := uinGroup.Group("/operate")
	uinSetGroup.POST("/login", sinking_server.HandleFunc(handleUin.Login))             //登陆QQ
	uinSetGroup.POST("/setting", sinking_server.HandleFunc(handleUin.Setting))         //修改功能设置
	uinSetGroup.POST("/password", sinking_server.HandleFunc(handleUin.ChangePassword)) //修改挂机密码
	uinSetGroup.POST("/server", sinking_server.HandleFunc(handleUin.ChangeServer))     //修改挂机地区

	//挂机节点
	handleServer := &uin.ControllerServer{}
	uinGroup.GET("/servers", sinking_server.HandleFunc(handleServer.List)) //地区列表
	handleLevel := &uin.ControllerLevel{}
	uinGroup.GET("/levels", sinking_server.HandleFunc(handleLevel.List)) //折扣列表
	handleOrder := &uin.ControllerOrder{}
	uinGroup.GET("/orders", sinking_server.HandleFunc(handleOrder.List)) //订单列表
	handleLog := &uin.ControllerLog{}
	uinGroup.GET("/logs", sinking_server.HandleFunc(handleLog.List)) //日志列表
}

// 后台系统路由
func loadAdminRoute(s *sinking_web.Engine) {
	g := s.Group("/admin")
	g.Use(middleware.AuthAdmin(), middleware2.Rbac())
	g.GET("/test", sinking_server.HandleFunc(admin.Index))
}
