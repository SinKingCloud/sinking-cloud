package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/controller/user/app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/controller/user/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/controller/user/oauth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/controller/user/ocr"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/http/controller/user/uin"
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

	//授权API
	handleAuth := &auth.ControllerAuth{}
	authGroup := g.Group("/auth")
	authGroup.Use(middleware2.ApiCheck())
	authGroup.GET("/check", sinking_server.HandleFunc(handleAuth.Check))     //授权检测
	authGroup.POST("/log", sinking_server.HandleFunc(handleAuth.Log))        //授权记录
	authGroup.GET("/upgrade", sinking_server.HandleFunc(handleAuth.Upgrade)) //升级信息
	authGroup.GET("/package", sinking_server.HandleFunc(handleAuth.Package)) //获取程序

	//应用API
	appGroup := g.Group("/app")
	appGroup.Use(middleware2.ApiCheck())
	handleApp := &app.ControllerApp{}
	appGroup.GET("/info", sinking_server.HandleFunc(handleApp.Info)) //应用信息

	//云端挂机API
	uinGroup := g.Group("/uin")
	uinGroup.Use(middleware2.ApiCheck())
	handleUin := &uin.ControllerUin{}
	uinGroup.GET("/self", sinking_server.HandleFunc(handleUin.Self))                    //账户信息
	uinGroup.POST("/buy", sinking_server.HandleFunc(handleUin.Buy))                     //业务下单
	uinGroup.GET("/servers", sinking_server.HandleFunc((&uin.ControllerServer{}).List)) //地区列表
	uinGroup.GET("/project", sinking_server.HandleFunc(handleUin.Project))              //功能列表
	uinGroup.GET("/activity", sinking_server.HandleFunc(handleUin.Activity))            //活动列表
	uinGroup.POST("/test", sinking_server.HandleFunc(handleUin.Test))                   //测试加速
	uinGroup.POST("/experience", sinking_server.HandleFunc(handleUin.Experience))       //体验套装
	uinSetGroup := uinGroup.Group("/operate")
	uinSetGroup.POST("/info", sinking_server.HandleFunc(handleUin.Info))               //挂机QQ详情
	uinSetGroup.POST("/login", sinking_server.HandleFunc(handleUin.Login))             //登陆挂机QQ
	uinSetGroup.POST("/qr_login", sinking_server.HandleFunc(handleUin.QrLogin))        //扫码验证
	uinSetGroup.POST("/sms", sinking_server.HandleFunc(handleUin.Sms))                 //短信验证
	uinSetGroup.POST("/delete", sinking_server.HandleFunc(handleUin.Delete))           //删除挂机QQ
	uinSetGroup.POST("/run", sinking_server.HandleFunc(handleUin.Job))                 //执行挂机任务
	uinSetGroup.POST("/setting", sinking_server.HandleFunc(handleUin.Setting))         //修改功能设置
	uinSetGroup.POST("/password", sinking_server.HandleFunc(handleUin.ChangePassword)) //修改挂机密码
	uinSetGroup.POST("/server", sinking_server.HandleFunc(handleUin.ChangeServer))     //修改挂机地区
	uinSetGroup.POST("/timing", sinking_server.HandleFunc(handleUin.ChangeTiming))     //修改挂机时间
	uinSetGroup.POST("/status", sinking_server.HandleFunc(handleUin.ChangeStatus))     //修改挂机状态
	uinToolGroup := uinGroup.Group("/tool")
	uinToolGroup.POST("/info", sinking_server.HandleFunc(handleUin.LevelInfo)) //获取QQ信息

	//全站登陆API
	oauthGroup := g.Group("/oauth")
	oauthGroup.Use(middleware2.ApiCheck())
	oauthProxy := &oauth.ControllerOauth{}
	oauthGroup.POST("/gen_ticket", sinking_server.HandleFunc(oauthProxy.GenTicket))      //生成ticket
	oauthGroup.POST("/get_user_info", sinking_server.HandleFunc(oauthProxy.GetUserInfo)) //获取授权用户信息
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

	//挂机应用
	uinGroup := g.Group("/uin")
	handleUin := &uin.ControllerUin{}
	uinGroup.GET("/list", sinking_server.HandleFunc(handleUin.List))              //挂机列表
	uinGroup.GET("/self", sinking_server.HandleFunc(handleUin.Self))              //账户信息
	uinGroup.GET("/count", sinking_server.HandleFunc(handleUin.Count))            //统计信息
	uinGroup.GET("/project", sinking_server.HandleFunc(handleUin.Project))        //功能列表
	uinGroup.GET("/activity", sinking_server.HandleFunc(handleUin.Activity))      //活动列表
	uinGroup.POST("/buy", sinking_server.HandleFunc(handleUin.Buy))               //业务下单
	uinGroup.POST("/test", sinking_server.HandleFunc(handleUin.Test))             //测试加速
	uinGroup.POST("/experience", sinking_server.HandleFunc(handleUin.Experience)) //体验套装
	uinSetGroup := uinGroup.Group("/operate")
	uinSetGroup.POST("/info", sinking_server.HandleFunc(handleUin.Info))               //挂机QQ详情
	uinSetGroup.POST("/login", sinking_server.HandleFunc(handleUin.Login))             //登陆挂机QQ
	uinSetGroup.POST("/qr_login", sinking_server.HandleFunc(handleUin.QrLogin))        //扫码验证
	uinSetGroup.POST("/sms", sinking_server.HandleFunc(handleUin.Sms))                 //短信验证
	uinSetGroup.POST("/delete", sinking_server.HandleFunc(handleUin.Delete))           //删除挂机QQ
	uinSetGroup.POST("/run", sinking_server.HandleFunc(handleUin.Job))                 //执行挂机任务
	uinSetGroup.POST("/setting", sinking_server.HandleFunc(handleUin.Setting))         //修改功能设置
	uinSetGroup.POST("/password", sinking_server.HandleFunc(handleUin.ChangePassword)) //修改挂机密码
	uinSetGroup.POST("/server", sinking_server.HandleFunc(handleUin.ChangeServer))     //修改挂机地区
	uinSetGroup.POST("/timing", sinking_server.HandleFunc(handleUin.ChangeTiming))     //修改挂机时间
	uinSetGroup.POST("/status", sinking_server.HandleFunc(handleUin.ChangeStatus))     //修改挂机状态
	uinToolGroup := uinGroup.Group("/tool")
	uinToolGroup.POST("/info", sinking_server.HandleFunc(handleUin.LevelInfo)) //获取QQ信息
	handleServer := &uin.ControllerServer{}
	uinGroup.GET("/servers", sinking_server.HandleFunc(handleServer.List)) //地区列表
	handleLevel := &uin.ControllerLevel{}
	uinGroup.GET("/levels", sinking_server.HandleFunc(handleLevel.List)) //折扣列表
	uinGroup.POST("/shop", sinking_server.HandleFunc(handleLevel.Buy))   //开通等级
	handleOrder := &uin.ControllerOrder{}
	uinGroup.GET("/orders", sinking_server.HandleFunc(handleOrder.List)) //订单列表
	handleLog := &uin.ControllerLog{}
	uinGroup.GET("/logs", sinking_server.HandleFunc(handleLog.List)) //日志列表

	//oauth应用
	handleOauth := &oauth.ControllerOauth{}
	oauthGroup := g.Group("/oauth")
	oauthWebGroup := oauthGroup.Group("/web")
	oauthWebGroup.GET("/list", sinking_server.HandleFunc(handleOauth.List))      //接入列表
	oauthWebGroup.POST("/delete", sinking_server.HandleFunc(handleOauth.Delete)) //删除接入
	oauthWebGroup.POST("/update", sinking_server.HandleFunc(handleOauth.Update)) //更新接入
	oauthWebGroup.POST("/create", sinking_server.HandleFunc(handleOauth.Create)) //新增接入

	//oauth应用
	handleOcr := &ocr.ControllerOcr{}
	ocrGroup := g.Group("/ocr")
	ocrGroup.POST("/test", sinking_server.HandleFunc(handleOcr.Test)) //测试接口

	//授权应用
	handleAuth := &auth.ControllerAuth{}
	authGroup := g.Group("/auth")
	authGroup.GET("/list", sinking_server.HandleFunc(handleAuth.List))              //授权列表
	authGroup.POST("/download", sinking_server.HandleFunc(handleAuth.Download))     //程序下载
	authGroup.GET("/app/info", sinking_server.HandleFunc(handleAuth.AppInfo))       //程序详情
	authGroup.GET("/app/version", sinking_server.HandleFunc(handleAuth.AppVersion)) //程序版本
	authGroup.GET("/shop/list", sinking_server.HandleFunc(handleAuth.ShopList))     //程序列表
	authGroup.POST("/shop/buy", sinking_server.HandleFunc(handleAuth.ShopBuy))      //购买程序
}

// 后台系统路由
func loadAdminRoute(s *sinking_web.Engine) {
	g := s.Group("/admin")
	g.Use(middleware.AuthAdmin(), middleware2.Rbac())
}
