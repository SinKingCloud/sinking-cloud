package balance

import (
	"service-pay/app/repository/balance"
	"service-pay/app/repository/balance_log"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/cache"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
)

// Service service接口
type Service interface {
	FindByUserId(userId int64, refresh bool) (*Balance, error)
	OperateMoney(userId int64, money float64, logType int, describe string) error                                                                         //操作用户余额并写入明细
	OperateMoneyCount(userId int64, money float64) error                                                                                                  //操作累计充值或消费金额
	SelectLog(where *balance_log.SelectBalanceLog, orderByField string, orderByType string, page int, pageSize int) (list []*Log, total int64, err error) //查询用户余额明细
}

// Session服务
type service struct {
	balanceRepository    balance.Interface
	balanceLogRepository balance_log.Interface
	database             *database.Database
	cache                cache.Interface
}

// NewService 实例化service
func NewService(balanceRepository balance.Interface, balanceLogRepository balance_log.Interface, database *database.Database, cache cache.Interface) *service {
	return &service{
		balanceRepository:    balanceRepository,
		balanceLogRepository: balanceLogRepository,
		database:             database,
		cache:                cache,
	}
}
