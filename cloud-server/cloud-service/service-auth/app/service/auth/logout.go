package auth

import "errors"

// Logout 注销登录
func (s *service) Logout(isAdmin bool, token string) error {
	if token == "" {
		return errors.New("您的登录状态异常")
	}
	if !isAdmin {
		return s.authService.WebUserLogout(token)
	}
	return s.authService.SystemUserLogout(token)
}
