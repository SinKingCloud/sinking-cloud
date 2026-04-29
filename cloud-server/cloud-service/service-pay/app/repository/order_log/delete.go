package order_log

import (
	"service-pay/app/model"

	"gorm.io/gorm"
)

// DeleteByOrderIds 根据订单ID批量删除数据
func (r *Repository) DeleteByOrderIds(orderIds []int64, tx ...*gorm.DB) error {
	if len(orderIds) <= 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(orderIds, 1000, func(batch interface{}) error {
			return db.Where("order_id IN ?", batch.([]int64)).Delete(&model.OrderLog{}).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
