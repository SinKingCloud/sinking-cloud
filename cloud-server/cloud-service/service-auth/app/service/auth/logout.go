package auth

// Logout 注销登录
func (s *service) Logout(isAdmin bool, token string) error {
	if !isAdmin {
		return s.authService.WebUserLogout(token)
	}
	return s.authService.SystemUserLogout(token)
}
