package message

import (
	"service-user/app/model"

	"gorm.io/gorm"
)

// DeleteByIdsAndUserId 根据消息ID集合和用户ID删除消息
func (r *Repository) DeleteByIdsAndUserId(ids []int64, userId int64, tx ...*gorm.DB) error {
	if len(ids) == 0 || userId <= 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
			return db.
				Where("id IN ?", batch.([]int64)).
				Where("user_id = ?", userId).
				Delete(&model.Message{}).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}

// DeleteByUserId 根据用户ID删除全部消息
func (r *Repository) DeleteByUserId(userId int64, tx ...*gorm.DB) error {
	if userId <= 0 {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Where("user_id = ?", userId).Delete(&model.Message{}).Error
}
