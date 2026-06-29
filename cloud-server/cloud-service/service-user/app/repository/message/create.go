package message

import (
	"service-user/app/model"

	"gorm.io/gorm"
)

// Create 创建消息
func (r *Repository) Create(message *model.Message) error {
	return r.Database.Db.Create(message).Error
}

// Creates 批量创建消息
func (r *Repository) Creates(messages []*model.Message, tx ...*gorm.DB) error {
	if len(messages) == 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(messages, 1000, func(batch interface{}) error {
			return db.Create(batch.([]*model.Message)).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
