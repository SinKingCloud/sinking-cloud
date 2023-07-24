package route

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-gateway/app/constant"
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
			groupName := tempV1["group"]   //服务组
			name := tempV1["name"]         //服务名
			protocol := tempV1["protocol"] //使用协议
			login := tempV1["login"]       //验证状态
			mode := 0                      //轮询算法0:随机/1:轮询/2:单机
			if tempV1["mode"] != nil {
				mode = int(tempV1["mode"].(float64))
			}
			tempGroup := route.Group(strings.Replace(k, "/*", "", -1))
			tempGroup2 := tempGroup.Group(strings.Replace(k1, "/*", "", -1))
			if login == "user" {
				tempGroup2.Use(middleware.AuthUser(), middleware2.CheckWebUser())
			}
			if login == "admin" {
				tempGroup2.Use(middleware.AuthAdmin(), middleware2.CheckSystemUser())
			}
			prefix := "http://"
			if protocol == "ws://" {
				prefix = protocol.(string)
			}
			patten := k
			tempGroup2.ANY(strings.Replace(k1, strings.Replace(k1, "/*", "", -1), "", -1), sinking_server.HandleFunc(func(s *sinking_server.Context) {
				s.Set(constant.ServerGroupName, groupName.(string))
				s.Set(constant.ServerName, name.(string))
				url, exists := util.Consul.GetService(groupName.(string), name.(string), mode)
				if exists {
					s.Set(constant.ServerAddr, url.Addr)
					err := s.Proxy(prefix+url.Addr, func(r *http.Request) *http.Request {
						r.URL.Path = strings.Replace(s.Request.URL.Path, patten, "", 1)
						r.URL.RawPath, r.RequestURI = r.URL.Path, r.URL.Path
						r.Header.Set(util.Conf.GetString("register.tokenKey"), util.Conf.GetString("register.tokenValue"))
						r.Header.Set(util2.RequestId, s.GetRequestId())
						ip := s.ClientIP(true)
						r.Header.Set("X-Forwarded-For", ip)
						r.Header.Set("X-Real-Ip", ip)
						return r
					})
					if err != nil {
						s.ErrorWithData("调用服务失败", err.Error())
					}
				} else {
					s.ErrorWithData("未发现有效服务", tempV1)
				}
			}))
		}
	}
}
