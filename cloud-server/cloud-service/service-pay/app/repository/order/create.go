package order

import "service-pay/app/model"

// Create 创建数据
func (r *Repository) Create(data *model.Order) error {
	if data == nil {
		return nil
	}
	return r.Database.Db.Create(data).Error
}
