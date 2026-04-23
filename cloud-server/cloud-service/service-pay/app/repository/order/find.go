package order

import (
	"errors"
	"service-pay/app/model"

	"gorm.io/gorm"
)

// FindByTradeNo 根据订单号查询记录
func (r *Repository) FindByTradeNo(tradeNo string) (*model.Order, error) {
	var data *model.Order
	err := r.Database.Db.Where("trade_no = ?", tradeNo).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}

// FindById 根据ID查询记录
func (r *Repository) FindById(id int64) (*model.Order, error) {
	var data *model.Order
	err := r.Database.Db.Where("id = ?", id).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil || data.Id <= 0 {
		return nil, nil
	}
	return data, nil
}
