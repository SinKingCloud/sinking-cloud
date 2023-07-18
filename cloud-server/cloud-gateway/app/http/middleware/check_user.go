package middleware

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
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
		webUser, err := service.GetWebUserServiceIns().FindByIdWithCache(jwtUser.Id, false)
		if err != nil {
			c.TokenError("登陆超时,请重新登陆", nil)
			c.Abort()
			return
		}
		webToken := jwt.GetLoginToken(webUser.LoginToken)
		//判断用户token一致性
		if webToken[types] == "" {
			c.TokenError("您的账户已注销登陆,请重新登陆", nil)
			c.Abort()
			return
		}
		if jwtToken[types] != webToken[types] && util2.Conf.GetString(config.ConfProjectEnv) != "dev" {
			c.TokenError("您的账户已在其他设备登陆,请重新登陆", nil)
			c.Abort()
			return
		}
		if webUser.Status != 0 {
			c.TokenError("您的账户已被禁止登陆,请联系管理员", nil)
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
		systemUser, err := service.GetSystemUserServiceIns().FindByIdWithCache(jwtUser.Id, false)
		if err != nil {
			c.TokenError("登陆超时,请重新登陆", nil)
			c.Abort()
			return
		}
		systemToken := jwt.GetLoginToken(systemUser.LoginToken)
		//判断用户token一致性
		if systemToken[types] == "" {
			c.TokenError("您的账户已注销登陆,请重新登陆", nil)
			c.Abort()
			return
		}
		if jwtToken[types] != systemToken[types] {
			c.TokenError("您的账户已在其他设备登陆,请重新登陆", nil)
			c.Abort()
			return
		}
		if systemUser.Status != 0 {
			c.TokenError("您的账户已被禁止登陆,请联系管理员", nil)
			c.Abort()
			return
		}
		c.Next()
	})
}
