package order

import (
	"service-pay/app/model"
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UpdateById 根据 id 更新
func (r *Repository) UpdateById(id int64, data *UpdateOrder, tx ...*gorm.DB) error {
	if id <= 0 || data == nil {
		return nil
	}
	updates := make(map[string]interface{})
	if data.Status != nil {
		updates["status"] = data.Status
	}
	if data.OutTradeNo != nil {
		updates["out_trade_no"] = data.OutTradeNo
	}
	if data.PayAmount != nil {
		updates["pay_amount"] = data.PayAmount
	}
	if data.RefundAmount != nil {
		updates["refund_amount"] = data.RefundAmount
	}
	if data.PayStatus != nil {
		updates["pay_status"] = data.PayStatus
	}
	if data.PayType != nil {
		updates["pay_type"] = data.PayType
	}
	if data.PayTime != nil {
		updates["pay_time"] = data.PayTime
	}
	if data.CompleteTime != nil {
		updates["complete_time"] = data.CompleteTime
	}
	if len(updates) == 0 {
		return nil
	}
	updates["update_time"] = str.DateTime(time.Now())
	db := r.Database.Db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Model(&model.Order{}).Where("id = ?", id).Updates(updates).Error
}
