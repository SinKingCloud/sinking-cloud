package user

import (
	"service-admin/app/model"
	queueUserLog "service-admin/app/queue/user_log"
	configRepository "service-admin/app/repository/config"
	sessionRepository "service-admin/app/repository/session"
	"service-admin/app/repository/user"
	userLogRepository "service-admin/app/repository/user_log"
	userRoleRepository "service-admin/app/repository/user_role"
	"service-admin/app/service/rbac"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/page"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

// Service service接口
type Service interface {
	GetInfo(userId int64, refresh bool) (user *model.User, err error)                                                    //获取用户信息
	Select(where *user.SelectUser, queryPage *page.Query) (*page.Result[*User], error)                                   //查询用户列表
	SelectLog(where *userLogRepository.SelectUserLog, queryPage *page.Query) (*page.Result[*model.UserLog], error)       //查询用户操作日志
	CreateLog(id int64, userId int64, requestId string, requestIp string, types int, title string, content string) error //异步创建用户操作日志
	Create(data *model.User) error                                                                                       //创建用户
	UpdateById(id int64, data *user.UpdateUser) error                                                                    //更新用户信息
	DeleteById(id int64) error                                                                                           //删除用户
	AssignRole(userId int64, roleIds []int64) error                                                                      //分配用户角色
	SetUserVerifyStatus(userId int64, allow bool, minute int)                                                            //设置用户安全验证状态
	CheckUserVerifyStatus(userId int64, clear bool) bool                                                                 // 判断用户安全验证状态
}

// Session服务
type service struct {
	repositoryUser     user.Interface
	repositoryUserLog  userLogRepository.Interface
	repositorySession  sessionRepository.Interface
	repositoryUserRole userRoleRepository.Interface
	repositoryConfig   configRepository.Interface
	rbacService        rbac.Service
	database           *database.Database
	cache              cache.Interface
	logQueue           *queue.UnicastClient[*queueUserLog.Task]
}

// NewService 实例化service
func NewService(repositoryUser user.Interface, repositoryUserLog userLogRepository.Interface, repositorySession sessionRepository.Interface, repositoryUserRole userRoleRepository.Interface, repositoryConfig configRepository.Interface, rbacService rbac.Service, database *database.Database, cache cache.Interface, logQueue *queue.UnicastClient[*queueUserLog.Task]) *service {
	ins := &service{
		repositoryUser:     repositoryUser,
		repositoryUserLog:  repositoryUserLog,
		repositorySession:  repositorySession,
		repositoryUserRole: repositoryUserRole,
		repositoryConfig:   repositoryConfig,
		rbacService:        rbacService,
		database:           database,
		cache:              cache,
		logQueue:           logQueue,
	}
	if logQueue != nil {
		logQueue.SetConsumer(ins.consumeLogTask)
	}
	return ins
}
