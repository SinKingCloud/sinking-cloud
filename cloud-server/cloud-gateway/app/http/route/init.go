package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-gateway/app/http/middleware"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"time"
)

// loadCommonMiddleware 加载全局中间件
func loadCommonMiddleware(s *sinking_web.Engine) {
	s.Use(middleware.Cors(), middleware.Log())
}

// loadErrorHandle 设置错误回调
func loadErrorHandle(s *sinking_web.Engine) {
	//设置错误回调
	s.SetErrorHandle(&sinking_web.ErrorHandel{
		//资源不存在错误
		NotFound: func(c *sinking_web.Context) {
			c.JSON(404, sinking_web.H{"code": 404, "message": "请求资源不存在"})
		},
		//系统错误
		Fail: func(c *sinking_web.Context, code int, message string) {
			c.JSON(500, sinking_web.H{"code": code, "message": message})
		},
	})
}

// Init 初始化server
func Init() {
	//设置是否以debug模式运行
	sinking_web.SetDebugMode(util.Conf.GetString("project.env") == "dev")
	//设置读写超时时间
	sinking_web.SetTimeOut(
		util.Conf.GetDuration("application.timeout.read")*time.Second,
		util.Conf.GetDuration("application.timeout.write")*time.Second,
	)
	//实例化一个http server
	r := sinking_web.Default()
	//加载错误handle
	loadErrorHandle(r)
	//加载全局中间件
	loadCommonMiddleware(r)
	//加载debug
	loadDebug(r)
	//加载路由
	loadRoute(r)
	//启动http server
	err := r.Run("0.0.0.0:" + util.Conf.GetString(config.ConfApplicationPort))
	if err != nil {
		panic(err)
		return
	}
}
