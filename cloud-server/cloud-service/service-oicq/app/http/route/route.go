package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/http/controller/user/uin"
	middleware2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/http/middleware"
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
	uinGroup := s.Group("/api")

	//云端挂机API
	uinGroup.Use(middleware2.ApiCheck())
	handleUin := &uin.ControllerUin{}
	uinGroup.GET("/self", sinking_server.HandleFunc(handleUin.Self))                    //账户信息
	uinGroup.GET("/auth", sinking_server.HandleFunc(handleUin.Auth))                    //权限信息
	uinGroup.POST("/buy", sinking_server.HandleFunc(handleUin.Buy))                     //业务下单
	uinGroup.GET("/servers", sinking_server.HandleFunc((&uin.ControllerServer{}).List)) //地区列表
	uinGroup.GET("/project", sinking_server.HandleFunc(handleUin.Project))              //功能列表
	uinGroup.GET("/activity", sinking_server.HandleFunc(handleUin.Activity))            //活动列表
	uinGroup.POST("/test", sinking_server.HandleFunc(handleUin.Test))                   //测试加速
	uinGroup.POST("/experience", sinking_server.HandleFunc(handleUin.Experience))       //体验套装
	uinSetGroup := uinGroup.Group("/operate")
	uinSetGroup.POST("/info", sinking_server.HandleFunc(handleUin.Info))                    //挂机QQ详情
	uinSetGroup.POST("/reset", sinking_server.HandleFunc(handleUin.ReSetDevice))            //重置设备
	uinSetGroup.POST("/login", sinking_server.HandleFunc(handleUin.Login))                  //登陆挂机QQ
	uinSetGroup.POST("/qr_login", sinking_server.HandleFunc(handleUin.QrLogin))             //扫码验证
	uinSetGroup.POST("/sms", sinking_server.HandleFunc(handleUin.Sms))                      //短信验证
	uinSetGroup.POST("/delete", sinking_server.HandleFunc(handleUin.Delete))                //删除挂机QQ
	uinSetGroup.POST("/run", sinking_server.HandleFunc(handleUin.Job))                      //执行挂机任务
	uinSetGroup.POST("/setting", sinking_server.HandleFunc(handleUin.Setting))              //修改功能设置
	uinSetGroup.POST("/password", sinking_server.HandleFunc(handleUin.ChangePassword))      //修改挂机密码
	uinSetGroup.POST("/server", sinking_server.HandleFunc(handleUin.ChangeServer))          //修改挂机地区
	uinSetGroup.POST("/timing", sinking_server.HandleFunc(handleUin.ChangeTiming))          //修改挂机时间
	uinSetGroup.POST("/status", sinking_server.HandleFunc(handleUin.ChangeStatus))          //修改挂机状态
	uinSetGroup.POST("/heart", sinking_server.HandleFunc(handleUin.Heart))                  //发送心跳
	uinSetGroup.POST("/cookies", sinking_server.HandleFunc(handleUin.GetCookies))           //获取cookies信息
	uinSetGroup.POST("/allow_qr", sinking_server.HandleFunc(handleUin.AllowQrCode))         //授权登录二维码
	uinSetGroup.POST("/allow_device", sinking_server.HandleFunc(handleUin.AllowDeviceLock)) //授权登录设备锁
	uinToolGroup := uinGroup.Group("/tool")
	uinToolGroup.POST("/info", sinking_server.HandleFunc(handleUin.LevelInfo)) //获取QQ信息
}

// 用户系统路由
func loadUserRoute(s *sinking_web.Engine) {
	uinGroup := s.Group("/user")
	uinGroup.Use(middleware.AuthUser())

	handleUin := &uin.ControllerUin{}
	uinGroup.GET("/list", sinking_server.HandleFunc(handleUin.List))              //挂机列表
	uinGroup.GET("/self", sinking_server.HandleFunc(handleUin.Self))              //账户信息
	uinGroup.GET("/auth", sinking_server.HandleFunc(handleUin.Auth))              //权限信息
	uinGroup.GET("/count", sinking_server.HandleFunc(handleUin.Count))            //统计信息
	uinGroup.GET("/project", sinking_server.HandleFunc(handleUin.Project))        //功能列表
	uinGroup.GET("/activity", sinking_server.HandleFunc(handleUin.Activity))      //活动列表
	uinGroup.POST("/buy", sinking_server.HandleFunc(handleUin.Buy))               //业务下单
	uinGroup.POST("/test", sinking_server.HandleFunc(handleUin.Test))             //测试加速
	uinGroup.POST("/experience", sinking_server.HandleFunc(handleUin.Experience)) //体验套装
	uinSetGroup := uinGroup.Group("/operate")
	uinSetGroup.POST("/info", sinking_server.HandleFunc(handleUin.Info))                    //挂机QQ详情
	uinSetGroup.POST("/reset", sinking_server.HandleFunc(handleUin.ReSetDevice))            //重置设备
	uinSetGroup.POST("/login", sinking_server.HandleFunc(handleUin.Login))                  //登陆挂机QQ
	uinSetGroup.POST("/qr_login", sinking_server.HandleFunc(handleUin.QrLogin))             //扫码验证
	uinSetGroup.POST("/sms", sinking_server.HandleFunc(handleUin.Sms))                      //短信验证
	uinSetGroup.POST("/delete", sinking_server.HandleFunc(handleUin.Delete))                //删除挂机QQ
	uinSetGroup.POST("/run", sinking_server.HandleFunc(handleUin.Job))                      //执行挂机任务
	uinSetGroup.POST("/setting", sinking_server.HandleFunc(handleUin.Setting))              //修改功能设置
	uinSetGroup.POST("/password", sinking_server.HandleFunc(handleUin.ChangePassword))      //修改挂机密码
	uinSetGroup.POST("/server", sinking_server.HandleFunc(handleUin.ChangeServer))          //修改挂机地区
	uinSetGroup.POST("/timing", sinking_server.HandleFunc(handleUin.ChangeTiming))          //修改挂机时间
	uinSetGroup.POST("/status", sinking_server.HandleFunc(handleUin.ChangeStatus))          //修改挂机状态
	uinSetGroup.POST("/heart", sinking_server.HandleFunc(handleUin.Heart))                  //发送心跳
	uinSetGroup.POST("/cookies", sinking_server.HandleFunc(handleUin.GetCookies))           //获取cookies信息
	uinSetGroup.POST("/allow_qr", sinking_server.HandleFunc(handleUin.AllowQrCode))         //授权登录二维码
	uinSetGroup.POST("/allow_device", sinking_server.HandleFunc(handleUin.AllowDeviceLock)) //授权登录设备锁
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
}

// 后台系统路由
func loadAdminRoute(s *sinking_web.Engine) {
	g := s.Group("/admin")
	g.Use(middleware.AuthAdmin(), middleware2.Rbac())
}
