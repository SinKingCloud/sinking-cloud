package route

import (
	"service-space/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadApi 加载公开API
func loadApi(s *sinking_web.Engine) {
	share := s.Group("/share")
	share.ANY("/info", context.HandleFunc(controller.Api.Share.Info)) //分享导图详情
}
