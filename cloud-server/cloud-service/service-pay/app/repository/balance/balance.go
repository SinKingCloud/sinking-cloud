package balance

import (
	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/database"
	"gorm.io/gorm"
)

// Interface 仓储接口
type Interface interface {
	FindOrCreateByUserId(userId int64) (*model.Balance, error)       // 通过用户ID查询余额，不存在则创建
	Money(userId int64, money float64, tx ...*gorm.DB) error         // 修改余额
	Frozen(userId int64, money float64, tx ...*gorm.DB) error        // 修改冻结金额
	TotalRecharge(userId int64, money float64, tx ...*gorm.DB) error // 修改累计充值金额
	TotalConsume(userId int64, money float64, tx ...*gorm.DB) error  // 修改累计消费金额
}

// Repository 仓储
type Repository struct {
	Database *database.Database
}

// NewRepository 创建仓储
func NewRepository(db *database.Database) *Repository {
	return &Repository{Database: db}
}
