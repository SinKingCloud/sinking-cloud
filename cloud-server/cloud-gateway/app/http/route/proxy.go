package route

import (
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	middleware2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-gateway/app/http/middleware"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
	"strings"
)

func InitProxyRouter(route *sinking_web.Engine) {
	conf := util.Consul.Config(group.CloudGateway).Name(config.Gateway).Viper()
	for k, v := range conf.AllSettings() {
		temp := v.(map[string]interface{})
		for k1, v1 := range temp {
			tempV1 := v1.(map[string]interface{})
			groupName := tempV1["group"]
			name := tempV1["name"]
			tempGroup := route.Group(k)
			if tempV1["login"] == "user" {
				tempGroup.Use(middleware.AuthUser(), middleware2.CheckWebUser())
			}
			if tempV1["login"] == "admin" {
				tempGroup.Use(middleware.AuthAdmin(), middleware2.CheckSystemUser())
			}
			patten := k
			tempGroup.ANY(k1, sinking_server.HandleFunc(func(s *sinking_server.Context) {
				url, exists := util.Consul.GetService(groupName.(string), name.(string))
				if exists {
					s.Proxy(patten, "http://"+url.Addr, func(r *http.Request) *http.Request {
						r.Host = s.Request.URL.Host
						r.URL.Path = strings.Replace(s.Request.URL.Path, patten, "", 1)
						r.Header.Set(util.Conf.GetString("register.tokenKey"), util.Conf.GetString("register.tokenValue"))
						r.Header.Set(util2.RequestId, s.GetRequestId())
						fmt.Println(r.Header.Clone())
						return r
					})
				} else {
					s.ErrorWithData("未发现有效服务", tempV1)
				}
			}))
		}
	}
}
