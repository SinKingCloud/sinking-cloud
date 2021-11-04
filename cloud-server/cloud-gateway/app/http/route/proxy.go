package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
)

func InitProxyRouter(route *sinking_web.Engine) {
	conf := util.Consul.Config(group.CloudGateway).Name(config.Gateway).Viper()
	for k, v := range conf.AllSettings() {
		value := v.(map[string]interface{})
		if value == nil || value["name"] == nil || value["group"] == nil {
			continue
		}
		groupName := value["group"]
		name := value["name"]
		route.ANY(k, sinking_server.HandleFunc(func(s *sinking_server.Context) {
			url, exists := util.Consul.GetService(groupName.(string), name.(string))
			if exists {
				s.Proxy(k, url.Addr, func(r *http.Request) *http.Request {
					return r
				})
			} else {
				s.ErrorWithData("未发现有效服务", value)
			}
		}))
	}
}
