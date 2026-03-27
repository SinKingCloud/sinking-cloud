package service

import (
	"service-user/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
)

type AuthService interface {
	Login(account string, password string, code string, userAgent string, ip string) (string, error) //账号登录
	Logout(token string) error                                                                       //登陆注销
	GetTokenInfo(token string) (*jwt.User, error)                                                    //获取登陆信息
	GetSessionInfo(token string) *model.Session                                                      //获取会话信息
}
