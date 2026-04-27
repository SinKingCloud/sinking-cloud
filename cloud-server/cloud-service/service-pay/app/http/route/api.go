package route

import (
	"service-pay/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApi 加载API
func loadApi(s *sinking_web.Engine) {
	s.ANY("/config", context.HandleFunc(controller.Api.Config))
	s.ANY("/submit", context.HandleFunc(controller.Api.Submit))
	s.ANY("/callback", context.HandleFunc(controller.Api.Callback))
}
