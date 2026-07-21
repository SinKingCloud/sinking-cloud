package mind_data

import (
	"service-space/app/model"

	"gorm.io/gorm"
)

// DeleteByIdAndUserId 根据ID和用户ID删除导图数据
func (r *Repository) DeleteByIdAndUserId(id int64, userId int64, tx ...*gorm.DB) error {
	if id <= 0 || userId <= 0 {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Where("id = ? AND user_id = ?", id, userId).Delete(&model.MindData{}).Error
}
