package order

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
	Id         int          `json:"id"`
	UserId     int          `json:"user_id"`
	PayType    int          `json:"pay_type"`
	OrderType  int          `json:"order_type"`
	Money      float64      `json:"money"`
	TradeNo    string       `json:"trade_no"`
	OutTradeNo string       `json:"out_trade_no"`
	Name       string       `json:"name"`
	Status     int          `json:"status"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}
