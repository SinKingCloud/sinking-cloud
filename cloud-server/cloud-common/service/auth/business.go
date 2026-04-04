package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/middleware"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rbac"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/context"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// CheckCaptcha 判断滑块验证码
func (*Service) CheckCaptcha(key string, x int, y int) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceAuth, "CheckCaptcha", &CheckCaptchaParam{
		Key: key,
		X:   x,
		Y:   y,
	}, &result)
	if err != nil {
		return false
	}
	return result
}

// CheckSms 判断短信验证码
func (*Service) CheckSms(phone string, code string) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceAuth, "CheckSms", &CheckSmsParam{
		Phone: phone,
		Code:  code,
	}, &result)
	if err != nil {
		return false
	}
	return result
}

// CheckMail 判断邮件验证码
func (*Service) CheckMail(mail string, code string) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceAuth, "CheckMail", &CheckMailParam{
		Mail: mail,
		Code: code,
	}, &result)
	if err != nil {
		return false
	}
	return result
}

// WebUserLogin 网站用户登陆
func (*Service) WebUserLogin(account string, password string, code string, userAgent string, ip string) (token string, err error) {
	err = rpc.InvokeFunc(constant.ServiceUser, "Login", &WebUserLoginParam{
		Account:   account,
		Password:  password,
		Code:      code,
		UserAgent: userAgent,
		Ip:        ip,
	}, &token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// SystemUserLogin 系统用户登陆
func (*Service) SystemUserLogin(account string, password string, code string, userAgent string, ip string) (token string, err error) {
	err = rpc.InvokeFunc(constant.ServiceAdmin, "Login", &SystemUserLoginParam{
		Account:   account,
		Password:  password,
		Code:      code,
		UserAgent: userAgent,
		Ip:        ip,
	}, &token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// WebUserLogout 网站用户注销登录
func (*Service) WebUserLogout(token string) error {
	err := rpc.InvokeFunc(constant.ServiceUser, "Logout", token, nil)
	if err != nil {
		return err
	}
	return nil
}

// SystemUserLogout 系统用户注销登录
func (*Service) SystemUserLogout(token string) error {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "Logout", token, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetAdminRouteGroup 获取后台路由group
func (*Service) GetAdminRouteGroup(s *sinking_web.Engine) *sinking_web.RouterGroup {
	g := s.Group("/" + constant.ServiceAdmin)
	g.Use(context.HandleFunc(middleware.CheckLogin))
	return g
}

// LoadAdminRouteWithRegister 加载后台路由组并注册
func (srv *Service) LoadAdminRouteWithRegister(g *sinking_web.RouterGroup, handle func(ANY func(path string, name string, h func(*context.Context)))) {
	var list []*rbac.RegisterPermissionsParam
	handle(func(path string, name string, h func(c *context.Context)) {
		if g != nil && h != nil {
			g.ANY(path, context.HandleFunc(h))
		}
		fullPath := strings.ReplaceAll(fmt.Sprintf("/%s/%s/%s", constant.ServiceAdmin, global.App.Consul.GetName(), path), "//", "/")
		list = append(list, &rbac.RegisterPermissionsParam{
			Name: name,
			Code: strings.Trim(strings.ReplaceAll(fullPath, "/", ":"), ":"),
			Path: fullPath,
		})
	})
	if list != nil && len(list) > 0 {
		go func(list []*rbac.RegisterPermissionsParam) {
			for i := 0; i < 60; i++ {
				result := rbac.GetIns().RegisterPermissions(list)
				if !result {
					time.Sleep(time.Second)
				} else {
					break
				}
			}
		}(list)
	}
}

// GetUserRouteGroup 获取用户路由group
func (*Service) GetUserRouteGroup(s *sinking_web.Engine) *sinking_web.RouterGroup {
	g := s.Group("/" + constant.ServiceUser)
	g.Use(context.HandleFunc(middleware.CheckLogin))
	return g
}
