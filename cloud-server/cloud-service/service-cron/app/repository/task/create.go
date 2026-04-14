package task

import "service-cron/app/model"

// Create 创建数据
func (r *Repository) Create(data *model.Task) error {
	if data == nil {
		return nil
	}
	return r.Database.Db.Create(data).Error
}
