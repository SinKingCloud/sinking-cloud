package balance

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// Balance 用户余额
type Balance struct {
	Money         float64      `gorm:"column:money" json:"money"`
	Frozen        float64      `gorm:"column:frozen" json:"frozen"`
	TotalRecharge float64      `gorm:"column:total_recharge" json:"total_recharge"`
	TotalConsume  float64      `gorm:"column:total_consume" json:"total_consume"`
	CreateTime    str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime    str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// Log 余额明细展示数据
type Log struct {
	Id           int64        `json:"id"`
	Type         int          `json:"type"`
	Money        float64      `json:"money"`
	MoneyBefore  float64      `json:"money_before"`
	MoneyAfter   float64      `json:"money_after"`
	FrozenBefore float64      `json:"frozen_before"`
	FrozenAfter  float64      `json:"frozen_after"`
	Describe     string       `json:"describe"`
	CreateTime   str.DateTime `json:"create_time"`
}
