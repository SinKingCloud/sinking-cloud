package message

import (
	"service-user/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateByIdAndUserId 根据消息ID和用户ID更新消息
func (r *Repository) UpdateByIdAndUserId(id int64, userId int64, data *UpdateMessage, tx ...*gorm.DB) error {
	if id <= 0 || userId <= 0 {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return r.update(db.Model(&model.Message{}).
		Where("id = ?", id).
		Where("user_id = ?", userId), data)
}

// UpdateByIdsAndUserId 根据消息ID列表和用户ID批量更新消息；ids 为空时更新该用户全部符合条件的消息
func (r *Repository) UpdateByIdsAndUserId(ids []int64, userId int64, data *UpdateMessage, tx ...*gorm.DB) error {
	if userId <= 0 {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	if len(ids) > 0 {
		execute := func(db *gorm.DB) error {
			return r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
				return r.update(db.Model(&model.Message{}).
					Where("id IN ?", batch.([]int64)).
					Where("user_id = ?", userId), data)
			})
		}
		if len(tx) > 0 && tx[0] != nil {
			return execute(tx[0])
		}
		return r.Database.Transaction(func(tx *gorm.DB) error {
			return execute(tx)
		})
	}
	return r.update(db.Model(&model.Message{}).Where("user_id = ?", userId), data)
}

// update 更新
func (r *Repository) update(db *gorm.DB, data *UpdateMessage) error {
	if data == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if data.Status != nil {
		updates["status"] = data.Status
	}
	if data.ReadTime != nil {
		updates["read_time"] = data.ReadTime
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	return db.Updates(updates).Error
}
