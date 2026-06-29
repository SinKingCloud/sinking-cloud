package message

import (
	"service-user/app/enum/message_status"
	"service-user/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateByIdAndUserId 根据消息ID和用户ID更新消息
func (r *Repository) UpdateByIdAndUserId(id int64, userId int64, data *UpdateMessage, tx ...*gorm.DB) error {
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
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Model(&model.Message{}).Where("id = ? AND user_id = ?", id, userId).Updates(updates).Error
}

// ReadByIdsAndUserId 根据消息ID列表和用户ID批量设置消息为已读
func (r *Repository) ReadByIdsAndUserId(ids []int64, userId int64, readTime str.DateTime, tx ...*gorm.DB) error {
	if len(ids) == 0 {
		return nil
	}
	updates := map[string]interface{}{
		"status":      message_status.Read,
		"read_time":   readTime,
		"update_time": readTime,
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
			return db.Model(&model.Message{}).
				Where("id IN ?", batch.([]int64)).
				Where("user_id = ?", userId).
				Where("status = ?", message_status.Unread).
				Updates(updates).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
