package auth

import (
	"service-user/app/model"
	"service-user/app/repository/session"
	"service-user/app/repository/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/jwt"
)

type Service interface {
	Login(account string, password string, code string, userAgent string, ip string) (string, error) //账号登录
	Logout(token string) error                                                                       //登陆注销
	GetTokenInfo(token string) (*jwt.User, error)                                                    //获取登陆信息
	GetSessionInfo(token string) *model.Session                                                      //获取会话信息
}

// Session服务
type service struct {
	authService       *auth.Service
	repositorySession session.Interface
	repositoryUser    user.Interface
	cache             cache.Interface
}

// NewService 实例化service
func NewService(authService *auth.Service, repositorySession session.Interface, repositoryUser user.Interface, cache cache.Interface) *service {
	return &service{
		authService:       authService,
		repositorySession: repositorySession,
		repositoryUser:    repositoryUser,
		cache:             cache,
	}
}
