package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/controller/pay"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/controller/person"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-user/app/http/controller/verify"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRoute 加载路由
func loadRoute(s *sinking_web.Engine) {
	loadPersonRoute(s)
	loadVerifyRoute(s)
	loadPayRoute(s)
}

// loadPersonRoute 加载个人操作路由
func loadPersonRoute(s *sinking_web.Engine) {
	r := s.Group("/person")
	handle := &person.ControllerPerson{}
	r.GET("/info", sinking_server.HandleFunc(handle.Info))                 //个人信息
	r.POST("/update/info", sinking_server.HandleFunc(handle.Update))       //更新资料
	r.POST("/update/password", sinking_server.HandleFunc(handle.Password)) //修改密码
	r.POST("/update/phone", sinking_server.HandleFunc(handle.Phone))       //修改手机
	r.GET("/auth/info", sinking_server.HandleFunc(handle.AuthInfo))        //实名信息
	r.POST("/auth/init", sinking_server.HandleFunc(handle.AuthInit))       //提交实名
}

// loadWebRoute 加载验证器路由
func loadVerifyRoute(s *sinking_web.Engine) {
	r := s.Group("/verify")
	handle := &verify.ControllerVerify{}
	r.POST("/sms", sinking_server.HandleFunc(handle.Sms))
}

// loadPayRoute 加载支付操作路由
func loadPayRoute(s *sinking_web.Engine) {
	r := s.Group("/pay")
	handle := &pay.ControllerPay{}
	r.GET("/order", sinking_server.HandleFunc(handle.OrderList))    //订单记录
	r.POST("/recharge", sinking_server.HandleFunc(handle.Recharge)) //充值
	r.POST("/log", sinking_server.HandleFunc(handle.LogList))       //消费记录
}
