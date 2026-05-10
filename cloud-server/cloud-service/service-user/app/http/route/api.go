package route

import (
	"service-user/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadUser 加载api
func loadApi(s *sinking_web.Engine) {
	g := s.Group("/api")

	callback := g.Group("/callback")
	{
		callback.ANY("/certificate", context.HandleFunc(controller.Api.Callback.Certificate)) //实名认证回调
	}
}
