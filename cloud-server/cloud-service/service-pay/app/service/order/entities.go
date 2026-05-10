package order

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// UserOrder 用户订单
type UserOrder struct {
	Id           int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	TradeNo      string       `gorm:"column:trade_no" json:"trade_no"`
	Name         string       `gorm:"column:name" json:"name"`
	Type         int          `gorm:"column:type" json:"type"`
	Amount       float64      `gorm:"column:amount" json:"amount"`
	Status       int          `gorm:"column:status" json:"status"`
	PayAmount    float64      `gorm:"column:pay_amount" json:"pay_amount"`
	RefundAmount float64      `gorm:"column:refund_amount" json:"refund_amount"`
	PayStatus    int          `gorm:"column:pay_status" json:"pay_status"`
	PayTime      str.DateTime `gorm:"column:pay_time" json:"pay_time"`
	ExpiredTime  str.DateTime `gorm:"column:expired_time" json:"expired_time"`
	CompleteTime str.DateTime `gorm:"column:complete_time" json:"complete_time"`
	CreateTime   str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime   str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// Extra 订单扩展信息
type Extra struct {
	RpcService string `json:"rpc_service"` //支付成功后调用的rpc服务名
	RpcFunc    string `json:"rpc_func"`    //支付成功后调用的rpc服务方法名
	Param      string `json:"param"`       //其他参数
}
