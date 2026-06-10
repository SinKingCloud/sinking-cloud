package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Order 订单表
type Order struct {
	Id           int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId       int64        `gorm:"column:user_id" json:"user_id"`
	ChannelId    int64        `gorm:"column:channel_id" json:"channel_id"`
	Name         string       `gorm:"column:name" json:"name"`
	Type         int          `gorm:"column:type" json:"type"`
	Extra        string       `gorm:"column:extra" json:"extra"`
	Amount       float64      `gorm:"column:amount" json:"amount"`
	Status       int          `gorm:"column:status" json:"status"`
	TradeNo      string       `gorm:"column:trade_no" json:"trade_no"`
	PayAmount    float64      `gorm:"column:pay_amount" json:"pay_amount"`
	RefundAmount float64      `gorm:"column:refund_amount" json:"refund_amount"`
	PayStatus    int          `gorm:"column:pay_status" json:"pay_status"`
	PayType      string       `gorm:"column:pay_type" json:"pay_type"`
	PayTime      str.DateTime `gorm:"column:pay_time" json:"pay_time"`
	ExpiredTime  str.DateTime `gorm:"column:expired_time" json:"expired_time"`
	CompleteTime str.DateTime `gorm:"column:complete_time" json:"complete_time"`
	CreateTime   str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime   str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Order) TableName() string {
	return "cloud_orders"
}

// BeforeCreate 创建前
func (t *Order) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Order) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
