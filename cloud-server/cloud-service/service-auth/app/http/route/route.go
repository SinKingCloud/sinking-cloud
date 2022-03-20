package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/controller/callback"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/controller/pay"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/controller/system"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/controller/upload"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/controller/verify"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-auth/app/http/controller/web"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	loadWebRoute(s)
	loadSystemRoute(s)
	loadVerifyRoute(s)
	loadUploadRoute(s)
	loadCallbackRoute(s)
	loadPayRoute(s)
}

// loadWebRoute 加载网站授权路由
func loadWebRoute(s *sinking_web.Engine) {
	r := s.Group("/web")
	handle := &web.ControllerWeb{}
	r.POST("/login", sinking_server.HandleFunc(handle.Login))
}

// loadWebRoute 加载后台授权路由
func loadSystemRoute(s *sinking_web.Engine) {
	r := s.Group("/system")
	handle := &system.ControllerSystem{}
	r.POST("/login", sinking_server.HandleFunc(handle.Login))
}

// loadWebRoute 加载验证器路由
func loadVerifyRoute(s *sinking_web.Engine) {
	r := s.Group("/verify")
	handle := &verify.ControllerVerify{}
	r.GET("/code", sinking_server.HandleFunc(handle.Code))
	r.POST("/sms", sinking_server.HandleFunc(handle.Sms))
}

// loadUploadRoute 加载上传路由
func loadUploadRoute(s *sinking_web.Engine) {
	r := s.Group("/upload")
	handle := &upload.ControllerUpload{}
	r.POST("/file", sinking_server.HandleFunc(handle.File))
}

// loadCallbackRoute 加载回调路由
func loadCallbackRoute(s *sinking_web.Engine) {
	r := s.Group("/callback")
	handle := &callback.ControllerCallback{}
	r.GET("/face", sinking_server.HandleFunc(handle.Face))
	r.POST("/alipay", sinking_server.HandleFunc(handle.Alipay))
}

// loadPayRoute 加载支付跳转路由
func loadPayRoute(s *sinking_web.Engine) {
	r := s.Group("/pay")
	handle := &pay.ControllerPay{}
	r.GET("/submit", sinking_server.HandleFunc(handle.Submit)) //跳转订单支付
	r.GET("/status", sinking_server.HandleFunc(handle.Status)) //查询订单状态
}
