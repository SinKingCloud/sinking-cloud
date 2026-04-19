package file

import (
	"service-driver/app/model"

	"gorm.io/gorm"
)

// DeleteByIds 根据ID集合删除数据
func (r *Repository) DeleteByIds(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return r.Database.BatchExecute(ids, 1000, func(batch interface{}) error {
			return tx.
				Where("id IN ?", batch.([]int64)).
				Delete(&model.File{}).Error
		})
	})
}
