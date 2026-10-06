package balance_log

import (
	"errors"
	"service-pay/app/model"

	"gorm.io/gorm"
)

// FindById 通过ID查询数据
func (r *Repository) FindById(id int64) (*model.BalanceLog, error) {
	if id <= 0 {
		return nil, nil
	}
	var data *model.BalanceLog
	err := r.Database.Db.Where("id = ?", id).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}
