package route

import (
	"service-socket/app/http/controller"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// loadRpc 加载rpc
func loadApi(s *sinking_web.Engine) {
	s.ANY("/", context.HandleFunc(controller.Socket))
}
