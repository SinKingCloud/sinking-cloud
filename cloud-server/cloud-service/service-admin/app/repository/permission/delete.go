package permission

import (
	"service-admin/app/model"

	"gorm.io/gorm"
)

// DeleteByIds 根据ID集合删除数据
func (r *Repository) DeleteByIds(ids []int64, tx ...*gorm.DB) error {
	if len(ids) == 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
			return db.Where("id IN ?", batch.([]int64)).Delete(&model.Permission{}).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
