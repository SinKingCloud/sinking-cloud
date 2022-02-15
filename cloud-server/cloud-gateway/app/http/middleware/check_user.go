package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// CheckWebUser 判断网站用户状态
func CheckWebUser() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		types := c.Request.Header.Get(util.JwtTokenNameType)
		if types == "" {
			c.TokenError("登陆超时,请重新登陆", nil)
			c.Abort()
			return
		}
		jwtUser := c.GetWebUserInfo()
		jwtToken := jwt.GetLoginToken(jwtUser.LoginToken)
		//获取用户最新信息
		webUser := service.GetWebUserServiceIns().FindByIdWithCache(jwtUser.Id, false)
		webToken := jwt.GetLoginToken(webUser.LoginToken)
		//判断用户token一致性
		if jwtToken[types] != webToken[types] {
			c.TokenError("您的账户已在其他设备登陆,请重新登陆", nil)
			c.Abort()
			return
		}
		c.Next()
	})
}

// CheckSystemUser 判断系统用户状态
func CheckSystemUser() sinking_web.HandlerFunc {
	return sinking_server.HandleFunc(func(c *sinking_server.Context) {
		types := c.Request.Header.Get(util.JwtTokenNameType)
		jwtUser := c.GetSystemUserInfo()
		jwtToken := jwt.GetLoginToken(jwtUser.LoginToken)
		//获取用户最新信息
		systemUser := service.GetSystemUserServiceIns().FindByIdWithCache(jwtUser.Id, false)
		systemToken := jwt.GetLoginToken(systemUser.LoginToken)
		//判断用户token一致性
		if jwtToken[types] != systemToken[types] {
			c.TokenError("您的账户已在其他设备登陆,请重新登陆", nil)
			c.Abort()
			return
		}
		c.Next()
	})
}
