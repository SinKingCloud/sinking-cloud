package channel

import "service-driver/app/model"

// Create 创建数据
func (r *Repository) Create(data *model.Channel) error {
	if data == nil {
		return nil
	}
	return r.Database.Db.Create(data).Error
}
