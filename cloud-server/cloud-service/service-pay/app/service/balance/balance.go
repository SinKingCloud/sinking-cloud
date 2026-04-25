package balance

import (
	"service-pay/app/repository/balance"
	"service-pay/app/repository/balance_log"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
)

// Service service接口
type Service interface {
	FindByUserId(userId int64, refresh bool) (*Balance, error)
	SelectLog(where *balance_log.SelectBalanceLog, orderByField string, orderByType string, page int, pageSize int) (list []*Log, total int64, err error) //查询用户余额明细
}

// Session服务
type service struct {
	balanceRepository    balance.Interface
	balanceLogRepository balance_log.Interface
	cache                cache.Interface
}

// NewService 实例化service
func NewService(balanceRepository balance.Interface, balanceLogRepository balance_log.Interface, cache cache.Interface) *service {
	return &service{
		balanceRepository:    balanceRepository,
		balanceLogRepository: balanceLogRepository,
		cache:                cache,
	}
}
