package scene

import "service-driver/app/model"

// Create 创建数据
func (r *Repository) Create(data *model.Scene) error {
	if data == nil {
		return nil
	}
	return r.Database.Db.Create(data).Error
}
