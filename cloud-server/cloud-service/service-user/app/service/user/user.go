package user

import (
	"service-user/app/model"
	"service-user/app/repository/real_auth"
	"service-user/app/repository/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Service service接口
type Service interface {
	GetInfo(userId int64, refresh bool) (user *model.User, err error)                                                            //获取用户信息
	UpdateById(id int64, data *user.UpdateUser) error                                                                            //更新用户信息
	GetRealAuthInfo(userId int64, refresh bool) (*RealAuthInfo, error)                                                           //获取实名信息
	InitRealAuth(userId int64, authType int, cardName, cardNo, companyName, companyNo, env string) (string, string, error)       //发起实名
	VerifyRealAuth(certifyId string) (*RealAuthInfo, error)                                                                      //查询实名状态
	SetUserVerifyStatus(userId int64, allow bool, minute int)                                                                    //设置用户安全验证状态
	CheckUserVerifyStatus(userId int64, clear bool) bool                                                                         //判断用户安全验证状态
	Select(where *user.SelectUser, orderByField string, orderByType string, page int, pageSize int) ([]*user.User, int64, error) //用户列表
}

// Session服务
type service struct {
	repositoryUser     user.Interface
	repositoryRealAuth real_auth.Interface
	config             *config.Service
	database           *database.Database
	cache              cache.Interface
}

// NewService 实例化service
func NewService(repositoryUser user.Interface, repositoryRealAuth real_auth.Interface, systemConfig *config.Service, database *database.Database, cache cache.Interface) *service {
	return &service{
		repositoryUser:     repositoryUser,
		repositoryRealAuth: repositoryRealAuth,
		config:             systemConfig,
		database:           database,
		cache:              cache,
	}
}
