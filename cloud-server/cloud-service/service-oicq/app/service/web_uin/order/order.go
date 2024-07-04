package order

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"sync"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

type Status int

const (
	Success Status = iota //订单成功
	Refund                //已退单
)

type Types int

const (
	All    Types = iota //全套
	Player              //挂机
	Sign                //签到
)

// AppUinOrder QQ订单记录
type AppUinOrder struct {
	TradeNo     string       `gorm:"column:trade_no" json:"trade_no"`
	UserId      int          `gorm:"column:user_id" json:"user_id"`
	AppId       int          `gorm:"column:app_id" json:"app_id"`
	UinId       int          `gorm:"column:uin_id" json:"uin_id"`
	Type        int          `gorm:"column:type" json:"type"`
	Money       float64      `gorm:"column:money" json:"money"`
	RefundMoney float64      `gorm:"column:refund_money" json:"refund_money"`
	Days        int          `gorm:"column:days" json:"days"`
	Name        string       `gorm:"column:name" json:"name"`
	Status      int          `gorm:"column:status" json:"status"`
	CreateTime  str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime  str.DateTime `gorm:"column:update_time" json:"update_time"`
}
