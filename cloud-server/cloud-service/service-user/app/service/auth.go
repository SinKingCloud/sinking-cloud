package service

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"

type AuthService interface {
	Login(account string, password string, code string, userAgent string, ip string) (string, error) //账号登录
	Logout(token string) error                                                                       //登陆注销
	GetTokenInfo(token string) (*jwt.User, error)                                                    //获取登陆信息
}
