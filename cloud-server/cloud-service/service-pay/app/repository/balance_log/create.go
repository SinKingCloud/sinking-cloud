package balance_log

import (
	"service-pay/app/model"

	"gorm.io/gorm"
)

// Create 创建数据
func (r *Repository) Create(data *model.BalanceLog, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Create(data).Error
}
