package log

import (
	"service-cron/app/model"

	"gorm.io/gorm"
)

// DeleteByTaskIds 根据任务ID集合删除任务数据
func (r *Repository) DeleteByTaskIds(taskIds []int64, tx ...*gorm.DB) error {
	if len(taskIds) == 0 {
		return nil
	}
	execute := func(db *gorm.DB) error {
		return r.Database.BatchExecute(taskIds, 1000, func(batch interface{}) error {
			return db.
				Where("task_id IN ?", batch.([]int64)).
				Delete(&model.Log{}).Error
		})
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return r.Database.Transaction(func(tx *gorm.DB) error {
		return execute(tx)
	})
}
