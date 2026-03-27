package auth

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service 接口
type Service interface {
	GetCaptcha(key string) (map[string]interface{}, error)                                                         //获取滑块验证码
	CheckCaptcha(key string, x int, y int) bool                                                                    //判断滑块验证码
	CheckSms(phone string, code string) bool                                                                       //判断短信验证码
	SendSms(phone string) bool                                                                                     //发送短信验证码
	CheckMail(mail string, code string) bool                                                                       //判断邮件验证码
	SendMail(mail string) bool                                                                                     //发送邮件验证码
	Login(isAdmin bool, account string, password string, code string, userAgent string, ip string) (string, error) //账号登录或注册
	Logout(isAdmin bool, token string) error                                                                       //账户注销登录
}

// Session服务
type service struct {
	cache         cache.Interface
	configService *config.Service
	authService   *auth.Service
}

// NewService 实例化service
func NewService(cache cache.Interface, configService *config.Service, authService *auth.Service) *service {
	return &service{
		cache:         cache,
		configService: configService,
		authService:   authService,
	}
}
