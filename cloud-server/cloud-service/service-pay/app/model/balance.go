package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Balance 用户余额
type Balance struct {
	Id            int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId        int64        `gorm:"column:user_id" json:"user_id"`
	Money         float64      `gorm:"column:money" json:"money"`
	Frozen        float64      `gorm:"column:frozen" json:"frozen"`
	TotalRecharge float64      `gorm:"column:total_recharge" json:"total_recharge"`
	TotalConsume  float64      `gorm:"column:total_consume" json:"total_consume"`
	CreateTime    str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime    str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Balance) TableName() string {
	return "cloud_balances"
}

// BeforeCreate 创建前
func (t *Balance) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Balance) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
