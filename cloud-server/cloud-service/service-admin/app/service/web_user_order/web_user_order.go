package web_user_order

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

// WebUserOrder 会员用户订单记录
type WebUserOrder struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	Account    string       `gorm:"column:account" json:"account"`
	NickName   string       `gorm:"column:nick_name" json:"nick_name"`
	Avatar     string       `gorm:"column:avatar" json:"avatar"`
	Email      string       `gorm:"column:email" json:"email"`
	Phone      int64        `gorm:"column:phone" json:"phone"`
	PayType    int          `gorm:"column:pay_type" json:"pay_type"`
	OrderType  int          `gorm:"column:order_type" json:"order_type"`
	Money      float64      `gorm:"column:money" json:"money"`
	TradeNo    string       `gorm:"column:trade_no" json:"trade_no"`
	OutTradeNo string       `gorm:"column:out_trade_no" json:"out_trade_no"`
	Name       string       `gorm:"column:name" json:"name"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}
