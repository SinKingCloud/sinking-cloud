package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// BalanceLog 用户余额记录表
type BalanceLog struct {
	Id           int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId       int64        `gorm:"column:user_id" json:"user_id"`
	Type         int          `gorm:"column:type" json:"type"`
	Money        float64      `gorm:"column:money" json:"money"`
	MoneyBefore  float64      `gorm:"column:money_before" json:"money_before"`
	MoneyAfter   float64      `gorm:"column:money_after" json:"money_after"`
	FrozenBefore float64      `gorm:"column:frozen_before" json:"frozen_before"`
	FrozenAfter  float64      `gorm:"column:frozen_after" json:"frozen_after"`
	Describe     string       `gorm:"column:describe" json:"describe"`
	CreateTime   str.DateTime `gorm:"column:create_time" json:"create_time"`
}

// TableName 获取表名
func (*BalanceLog) TableName() string {
	return "cloud_balance_logs"
}

// BeforeCreate 创建前
func (t *BalanceLog) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}
