package proxy

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_plugin/plugin"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_plugin/user_plugin"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

func FilterPlugins(isAdmin bool, s *sinking_web.RouterGroup) {
	s.ANY("/*", sinking_server.HandleFunc(func(c *sinking_server.Context) {
		//匹配应用
		temps := strings.Split(strings.Trim(c.Request.RequestURI, "/"), "/")
		if len(temps) < 2 {
			c.Error("请求错误,该请求路径不存在")
			c.Abort()
			return
		}
		app, err := service.Plugin.FindByUrlWithCache(temps[1], false)
		if err != nil || app == nil {
			c.Error("该应用不存在,请咨询管理员")
			c.Abort()
			return
		}
		user := c.GetWebUserInfo()
		if !isAdmin {
			if app.Status != int(plugin.StatusNormal) {
				c.Error("该应用已暂停提供服务,详情请咨询管理员")
				c.Abort()
				return
			}
			if user == nil {
				c.Error("获取登录用户信息失败,详情请咨询管理员")
				c.Abort()
				return
			}
			//付费插件判断到期时间
			if app.IsFree == int(plugin.IsFreeNo) {
				//查询用户插件信息
				userApp, err := service.UserPlugin.FindByUserIdAndPluginIdWithCache(user.Id, app.Id, false)
				if err != nil || userApp == nil {
					c.NotAllow("您的账号未开通该应用,请先开通应用", nil)
					c.Abort()
					return
				}
				if userApp.Status != int(user_plugin.StatusNormal) {
					c.Error("您的应用已被禁止提供服务,详情请咨询管理员")
					c.Abort()
					return
				}
				stamp, err := time.ParseInLocation("2006-01-02 15:04:05", userApp.ExpireTime.String(), time.Local)
				if err != nil || stamp.Unix() <= time.Now().Unix() {
					c.NotAllow("您的应用已到期,请续费应用", nil)
					c.Abort()
					return
				}
			}
		}
		ip := c.ClientIP(true)
		//代理应用请求(ip_hash模式)
		hash := ip
		if user != nil {
			hash = hash + "_" + strconv.Itoa(user.Id)
		}
		url, exists := util.Consul.GetServiceAddrByHash(app.ServiceGroup, app.ServiceName, hash)
		if exists {
			err = c.Proxy("http://"+url, log.New(io.Discard, "", 0), func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy) {
				r.URL.Path = strings.Replace(c.Request.URL.Path, "/"+app.FilterUrl, "", 1)
				r.URL.RawPath, r.RequestURI = r.URL.Path, r.URL.Path
				r.Header.Set(util.Conf.GetString(config.ConfRegisterTokenKey), util.Conf.GetString(config.ConfRegisterTokenValue))
				r.Header.Set(util2.RequestId, c.GetRequestId())
				r.Header.Set("X-Forwarded-For", ip)
				r.Header.Set("X-Real-Ip", ip)
			}, func(w http.ResponseWriter, r *http.Request, err error) {
				c.ErrorWithData("请求错误,服务已下线", nil)
			})
			if err != nil {
				c.ErrorWithData("请求错误,请联系管理员", err)
			}
		} else {
			c.Error("应用实例已下线,请联系管理员")
		}
	}))
}
