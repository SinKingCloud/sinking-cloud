package web_user_order

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

// WebUserOrder 会员用户订单记录
type WebUserOrder struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	Account    string       `json:"account"`
	NickName   string       `json:"nick_name"`
	Avatar     string       `json:"avatar"`
	Email      string       `json:"email"`
	Phone      int64        `json:"phone"`
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
