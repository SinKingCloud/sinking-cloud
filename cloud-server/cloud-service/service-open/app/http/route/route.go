package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/http/controller/admin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/http/controller/api"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/http/controller/user/app"
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
}

// 后台系统路由
func loadAdminRoute(s *sinking_web.Engine) {
	g := s.Group("/admin")
	g.Use(middleware.AuthAdmin(), middleware2.Rbac())
	g.GET("/test", sinking_server.HandleFunc(admin.Index))
}
