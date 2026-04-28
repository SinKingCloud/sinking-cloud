package route

import (
	"service-pay/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApi 加载API
func loadApi(s *sinking_web.Engine) {
	s.ANY("/config", context.HandleFunc(controller.Api.Config))                 //支付配置
	s.ANY("/submit", context.HandleFunc(controller.Api.Submit))                 //支付提交
	s.ANY("/callback/:id/:method", context.HandleFunc(controller.Api.Callback)) //支付回调
}
