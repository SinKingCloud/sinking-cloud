package user

import (
	"service-admin/app/model"
	sessionRepository "service-admin/app/repository/session"
	"service-admin/app/repository/user"
	userRoleRepository "service-admin/app/repository/user_role"
	configService "service-admin/app/service/config"
	"service-admin/app/service/rbac"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Service service接口
type Service interface {
	GetInfo(userId int64, refresh bool) (user *model.User, err error)                                                       //获取用户信息
	Select(where *user.SelectUser, orderByField string, orderByType string, page int, pageSize int) ([]*User, int64, error) //查询用户列表
	Create(data *model.User) error                                                                                          //创建用户
	UpdateById(id int64, data *user.UpdateUser) error                                                                       //更新用户信息
	DeleteById(id int64) error                                                                                              //删除用户
	AssignRole(userId int64, roleIds []int64) error                                                                         //分配用户角色
	SetUserVerifyStatus(userId int64, allow bool, minute int)                                                               //设置用户安全验证状态
	CheckUserVerifyStatus(userId int64, clear bool) bool                                                                    // 判断用户安全验证状态
}

// Session服务
type service struct {
	repositoryUser     user.Interface
	repositorySession  sessionRepository.Interface
	repositoryUserRole userRoleRepository.Interface
	configService      configService.Service
	rbacService        rbac.Service
	database           *database.Database
	cache              cache.Interface
}

// NewService 实例化service
func NewService(repositoryUser user.Interface, repositorySession sessionRepository.Interface, repositoryUserRole userRoleRepository.Interface, configService configService.Service, rbacService rbac.Service, database *database.Database, cache cache.Interface) *service {
	return &service{
		repositoryUser:     repositoryUser,
		repositorySession:  repositorySession,
		repositoryUserRole: repositoryUserRole,
		configService:      configService,
		rbacService:        rbacService,
		database:           database,
		cache:              cache,
	}
}
