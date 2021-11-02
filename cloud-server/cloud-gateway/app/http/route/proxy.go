package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

func InitProxyRouter(route *sinking_web.Engine) {
	route.GET("/", func(context *sinking_web.Context) {
		response.Success(context, "success")
	})
}
