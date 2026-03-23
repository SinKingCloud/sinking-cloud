package auth

// Login 账号登录
func (s *service) Login(isAdmin bool, account string, password string, code string, userAgent string, ip string) (string, error) {
	if !isAdmin {
		return s.authService.WebUserLogin(account, password, code, userAgent, ip)
	}
	return s.authService.SystemUserLogin(account, password, code, userAgent, ip)
}
