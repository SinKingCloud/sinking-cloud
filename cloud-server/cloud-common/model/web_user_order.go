package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// WebUserOrder 会员用户订单记录
type WebUserOrder struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	PayType    int          `gorm:"column:pay_type" json:"pay_type"`
	OrderType  int          `gorm:"column:order_type" json:"order_type"`
	Money      float64      `gorm:"column:money" json:"money"`
	TradeNo    string       `gorm:"column:trade_no" json:"trade_no"`
	OutTradeNo string       `gorm:"column:out_trade_no" json:"out_trade_no"`
	Name       string       `gorm:"column:name" json:"name"`
	Param      string       `gorm:"column:param" json:"param"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (WebUserOrder) TableName() string {
	return "web_user_orders"
}

// BeforeCreate 创建前
func (t *WebUserOrder) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *WebUserOrder) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
