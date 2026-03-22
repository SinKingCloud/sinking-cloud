package service

type AuthService interface {
	Login(account string, password string, code string, userAgent string, ip string) (string, error) //账号登录
}
