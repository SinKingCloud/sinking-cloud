package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/constant"
	app2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-open/app/service/web_user"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"strconv"
)

// ApiCheck api访问控制
func ApiCheck() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		temp := c.Request.Header.Get("app_id")
		appKey := c.Request.Header.Get("app_key")
		appId, err := strconv.Atoi(temp)
		if temp == "" || err != nil {
			c.TokenError("应用ID不能为空", nil)
			c.Abort()
			return
		}
		if appKey == "" {
			c.TokenError("应用密匙不能为空", nil)
			c.Abort()
			return
		}
		app, err := app2.GetIns().FindByIdWithCache(appId, false)
		if err != nil || app == nil || app.Status != 0 || app.Key != appKey {
			c.TokenError("获取应用状态异常,请检查应用状态及应用密匙", nil)
			c.Abort()
			return
		}
		u, err := web_user.GetIns().FindByIdWithCache(app.UserId, false)
		if err != nil || u == nil || u.Status != 0 || u.RealAuth != 0 {
			c.TokenError("获取应用所属账户异常,请检查账户状态以及是否实名", nil)
			c.Abort()
			return
		}
		c.Set(constant.AppInfo, app)
		c.SetWebUserInfo(&jwt.WebUser{
			Id:         u.Id,
			LoginToken: u.LoginToken,
			LoginIp:    u.LoginIp,
			LoginTime:  u.LoginTime,
			Status:     u.Status,
			RealAuth:   u.RealAuth,
		})
		c.Next()
	})
}
