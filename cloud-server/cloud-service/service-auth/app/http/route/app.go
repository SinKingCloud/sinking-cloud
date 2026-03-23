package route

import (
	"service-auth/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApp 加载app路由
func loadApp(s *sinking_web.Engine) {
	loadRpc(s)
	loadApi(s)
}

// loadApi
func loadApi(s *sinking_web.Engine) {
	s.ANY("/login", context.HandleFunc(controller.Login))     //账户登陆
	s.ANY("/logout", context.HandleFunc(controller.Logout))   //登陆注销
	s.ANY("/captcha", context.HandleFunc(controller.Captcha)) //滑块验证
	s.ANY("/sms", context.HandleFunc(controller.Sms))         //发送短信
	s.ANY("/mail", context.HandleFunc(controller.Mail))       //发送邮件
}
