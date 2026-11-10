package auth

import "errors"

// Logout 注销登录
func (s *service) Logout(isAdmin bool, token string) error {
	if s.authService == nil {
		return errors.New("认证服务未初始化")
	}
	if !isAdmin {
		return s.authService.WebUserLogout(token)
	}
	return s.authService.SystemUserLogout(token)
}
