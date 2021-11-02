package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
)

func InitProxyRouter(route *sinking_web.Engine) {
	conf := util.Consul.Config(group.CloudCommon).Name(config.Gateway).Viper()
	for k, v := range conf.AllSettings() {
		route.GET(k, func(s *sinking_web.Context) {
			url, exists := util.Consul.GetService("default", v.(string))
			if exists {
				s.Proxy(k, url.Addr, func(r *http.Request) *http.Request {
					return r
				})
			} else {
				response.Error(s, v)
			}
		})
	}
}
