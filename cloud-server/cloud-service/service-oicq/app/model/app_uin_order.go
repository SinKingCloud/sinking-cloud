package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppUinOrder QQ订单记录
type AppUinOrder struct {
	Id          int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId      int          `gorm:"column:user_id" json:"user_id"`
	AppId       int          `gorm:"column:app_id" json:"app_id"`
	UinId       int          `gorm:"column:uin_id" json:"uin_id"`
	Type        int          `gorm:"column:type" json:"type"`
	TradeNo     string       `gorm:"column:trade_no" json:"trade_no"`
	Money       float64      `gorm:"column:money" json:"money"`
	RefundMoney float64      `gorm:"column:refund_money" json:"refund_money"`
	Days        int          `gorm:"column:days" json:"days"`
	Name        string       `gorm:"column:name" json:"name"`
	Status      int          `gorm:"column:status" json:"status"`
	CreateTime  str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime  str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*AppUinOrder) TableName() string {
	return "app_uin_orders"
}

// BeforeCreate 创建前
func (t *AppUinOrder) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppUinOrder) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
