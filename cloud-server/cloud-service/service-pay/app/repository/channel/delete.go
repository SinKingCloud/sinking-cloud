package channel

import (
	"service-pay/app/model"
)

// DeleteById 根据ID删除数据
func (r *Repository) DeleteById(id int64) error {
	if id <= 0 {
		return nil
	}
	return r.Database.Db.Where("id = ?", id).
		Delete(&model.Channel{}).Error
}
