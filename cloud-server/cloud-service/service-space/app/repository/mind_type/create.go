package mind_type

import (
	"service-space/app/model"

	"gorm.io/gorm"
)

// Create 创建导图分类
func (r *Repository) Create(data *model.MindType, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Create(data).Error
}
