package channel

import (
	"errors"
	"service-driver/app/model"

	"gorm.io/gorm"
)

// FindById 根据ID查询记录
func (r *Repository) FindById(id int64) (*model.Channel, error) {
	var data *model.Channel
	err := r.Database.Db.Where("id = ?", id).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}
