package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/controller/proxy"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/controller/user/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/controller/user/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/controller/user/oauth"
	middleware2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/middleware"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	loadApiRoute(s)
	loadUserRoute(s)
	loadAdminRoute(s)
}

// 加载api路由
func loadApiRoute(s *sinking_web.Engine) {
	g := s.Group("/api")
	g.Use(middleware2.ApiCheck())
	//授权API
	handleAuth := &auth.ControllerAuth{}
	authGroup := g.Group("/auth")
	authGroup.GET("/check", sinking_server.HandleFunc(handleAuth.Check))     //授权检测
	authGroup.POST("/log", sinking_server.HandleFunc(handleAuth.Log))        //授权记录
	authGroup.GET("/upgrade", sinking_server.HandleFunc(handleAuth.Upgrade)) //升级信息
	authGroup.GET("/package", sinking_server.HandleFunc(handleAuth.Package)) //获取程序

	//应用API
	appGroup := g.Group("/app")
	handleApp := &app.ControllerApp{}
	appGroup.GET("/info", sinking_server.HandleFunc(handleApp.Info)) //应用信息

	//授权登录API
	oauthGroup := g.Group("/oauth")
	oauthProxy := &oauth.ControllerOauth{}
	oauthGroup.POST("/gen_ticket", sinking_server.HandleFunc(oauthProxy.GenTicket))      //生成ticket
	oauthGroup.POST("/get_user_info", sinking_server.HandleFunc(oauthProxy.GetUserInfo)) //获取授权用户信息

	//代理应用
	proxy.FilterPlugins(false, g)
}

// 用户系统路由
func loadUserRoute(s *sinking_web.Engine) {
	g := s.Group("/user")
	g.Use(middleware.AuthUser())
	//用户应用
	handle := &app.ControllerApp{}
	appGroup := g.Group("/app")
	appGroup.GET("/list", sinking_server.HandleFunc(handle.List))      //应用列表
	appGroup.POST("/create", sinking_server.HandleFunc(handle.Create)) //应用创建
	appGroup.POST("/update", sinking_server.HandleFunc(handle.Update)) //应用更新
	appGroup.POST("/delete", sinking_server.HandleFunc(handle.Delete)) //应用删除

	//授权应用
	handleAuth := &auth.ControllerAuth{}
	authGroup := g.Group("/auth")
	authGroup.GET("/list", sinking_server.HandleFunc(handleAuth.List))              //授权列表
	authGroup.POST("/download", sinking_server.HandleFunc(handleAuth.Download))     //程序下载
	authGroup.GET("/app/info", sinking_server.HandleFunc(handleAuth.AppInfo))       //程序详情
	authGroup.GET("/app/version", sinking_server.HandleFunc(handleAuth.AppVersion)) //程序版本
	authGroup.GET("/shop/list", sinking_server.HandleFunc(handleAuth.ShopList))     //程序列表
	authGroup.POST("/shop/buy", sinking_server.HandleFunc(handleAuth.ShopBuy))      //购买程序

	//授权登录API
	handleOauth := &oauth.ControllerOauth{}
	oauthWebGroup := g.Group("/oauth")
	oauthWebGroup.GET("/list", sinking_server.HandleFunc(handleOauth.List))      //接入列表
	oauthWebGroup.POST("/delete", sinking_server.HandleFunc(handleOauth.Delete)) //删除接入
	oauthWebGroup.POST("/update", sinking_server.HandleFunc(handleOauth.Update)) //更新接入
	oauthWebGroup.POST("/create", sinking_server.HandleFunc(handleOauth.Create)) //新增接入

	//代理应用
	proxy.FilterPlugins(false, g)
}

// 后台系统路由
func loadAdminRoute(s *sinking_web.Engine) {
	g := s.Group("/admin")
	g.Use(middleware.AuthAdmin(), middleware2.Rbac())

	//代理应用
	proxy.FilterPlugins(true, g)
}
