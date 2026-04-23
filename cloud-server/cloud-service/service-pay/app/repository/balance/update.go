package balance

import (
	"time"

	"service-pay/app/model"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Money 修改余额。
func (r *Repository) Money(userId int64, money float64, tx ...*gorm.DB) error {
	return r.updateAmount(userId, "money", money, tx...)
}

// Frozen 修改冻结金额。
func (r *Repository) Frozen(userId int64, money float64, tx ...*gorm.DB) error {
	return r.updateAmount(userId, "frozen", money, tx...)
}

// TotalRecharge 修改累计充值金额。
func (r *Repository) TotalRecharge(userId int64, money float64, tx ...*gorm.DB) error {
	return r.updateAmount(userId, "total_recharge", money, tx...)
}

// TotalConsume 修改累计消费金额。
func (r *Repository) TotalConsume(userId int64, money float64, tx ...*gorm.DB) error {
	return r.updateAmount(userId, "total_consume", money, tx...)
}

// updateAmount 修改余额相关金额字段。
func (r *Repository) updateAmount(userId int64, field string, money float64, tx ...*gorm.DB) error {
	if userId <= 0 || field == "" || money == 0 {
		return nil
	}
	operator := "+"
	value := money
	if money < 0 {
		operator = "-"
		value = -money
	}
	updates := map[string]interface{}{
		field:         gorm.Expr(field+" "+operator+" ?", value),
		"update_time": str.DateTime(time.Now()),
	}
	execute := func(db *gorm.DB) error {
		return db.
			Model(&model.Balance{}).
			Where("user_id = ?", userId).
			Updates(updates).Error
	}
	if len(tx) > 0 && tx[0] != nil {
		return execute(tx[0])
	}
	return execute(r.Database.Db)
}
