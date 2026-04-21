package task

import (
	"errors"
	"service-cron/app/model"

	"gorm.io/gorm"
)

// FindById 根据ID查询
func (r *Repository) FindById(id int64) (*model.Task, error) {
	var task *model.Task
	err := r.Database.Db.
		Model(&model.Task{}).
		Where("id = ?", id).
		First(&task).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if task == nil || task.Id <= 0 {
		return nil, nil
	}
	return task, nil
}
