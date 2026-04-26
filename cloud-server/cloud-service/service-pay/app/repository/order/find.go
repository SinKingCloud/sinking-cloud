package order

import (
	"errors"
	"service-pay/app/enum/order_pay_status"
	"service-pay/app/enum/order_status"
	"service-pay/app/model"
	"time"

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

// CheckChannelHasUnPayOrder 根据通道ID判断是否存在未过期且未支付订单
func (r *Repository) CheckChannelHasUnPayOrder(channelId int64) (bool, error) {
	if channelId <= 0 {
		return false, nil
	}
	var data *model.Order
	err := r.Database.Db.Model(&model.Order{}).
		Select("id").
		Where(
			"channel_id = ? AND status = ? AND pay_status = ? AND expired_time > ?",
			channelId,
			order_status.UnPay,
			order_pay_status.UnPay,
			time.Now().Format("2006-01-02 15:04:05"),
		).
		First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return data != nil && data.Id > 0, nil
}
