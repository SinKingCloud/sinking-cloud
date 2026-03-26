package global

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// checkLogin 判断是否登陆
func (a *Application) checkLogin(c *context.Context) {
	//登陆验证
	token := c.Request.Header.Get(constant.JwtTokenName)
	if token == "" {
		c.NotLogin("您还未登陆,请先登陆账户")
		c.Abort()
		return
	}
	key := jwt.CheckToken(token)
	if key == nil || key.User == nil || key.User.Id <= 0 || key.User.Token == "" {
		c.TokenError("登陆超时,请重新登陆")
		c.Abort()
		return
	}
	c.SetUserInfo(key.User)
	c.Next()
}

// GetAdminRouteGroup 获取后台路由组
func (a *Application) GetAdminRouteGroup(s *sinking_web.Engine) *sinking_web.RouterGroup {
	g := s.Group("/" + constant.ServiceAdmin)
	return g.Use(context.HandleFunc(a.checkLogin))
}

// GetUserRouteGroup 获取用户路由组
func (a *Application) GetUserRouteGroup(s *sinking_web.Engine) *sinking_web.RouterGroup {
	g := s.Group("/" + constant.ServiceUser)
	return g.Use(context.HandleFunc(a.checkLogin))
}
