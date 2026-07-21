package mind_type

import (
	"service-space/app/model"

	"gorm.io/gorm"
)

// CountByUserId 根据用户ID统计自定义导图分类数量
func (r *Repository) CountByUserId(userId int64, tx ...*gorm.DB) (int64, error) {
	if userId <= 0 {
		return 0, nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	var total int64
	err := db.Model(&model.MindType{}).Where("user_id = ?", userId).Count(&total).Error
	return total, err
}
