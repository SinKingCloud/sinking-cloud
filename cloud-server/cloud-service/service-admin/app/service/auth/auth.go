package auth

import (
	"service-admin/app/model"
	"service-admin/app/repository/session"
	"service-admin/app/repository/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
)

// Service 接口
type Service interface {
	Login(account string, password string, code string, userAgent string, ip string) (string, error)               //账号登录
	Logout(token string) error                                                                                     //登陆注销
	GetSessionInfo(token string) *model.Session                                                                    //获取会话信息
	SelectSessionByUserId(userId int64, where *session.SelectSession, page *page.Query) ([]*Session, int64, error) //会话列表
	LogoutByUserIdAndTokens(userId int64, tokens []string) error                                                   //批量注销会话
	ClearExpireSessions() error                                                                                    //清理过期会话
}

// Session服务
type service struct {
	authService       *auth.Service
	repositorySession session.Interface
	repositoryUser    user.Interface
	database          *database.Database
	cache             cache.Interface
}

// NewService 实例化service
func NewService(authService *auth.Service, repositorySession session.Interface, repositoryUser user.Interface, database *database.Database, cache cache.Interface) *service {
	return &service{
		authService:       authService,
		repositorySession: repositorySession,
		repositoryUser:    repositoryUser,
		database:          database,
		cache:             cache,
	}
}
