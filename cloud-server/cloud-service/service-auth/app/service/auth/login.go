package auth

import "errors"

// Login 账号登录
func (s *service) Login(isAdmin bool, account string, password string, code string, userAgent string, ip string) (string, error) {
	if s.authService == nil {
		return "", errors.New("认证服务未初始化")
	}
	if !isAdmin {
		return s.authService.WebUserLogin(account, password, code, userAgent, ip)
	}
	return s.authService.SystemUserLogin(account, password, code, userAgent, ip)
}
