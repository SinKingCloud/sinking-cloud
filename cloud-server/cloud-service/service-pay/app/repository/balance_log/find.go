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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}
