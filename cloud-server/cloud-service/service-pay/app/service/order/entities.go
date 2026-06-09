package order

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// UserOrder 用户订单
type UserOrder struct {
	Id           int64        `json:"id"`
	TradeNo      string       `json:"trade_no"`
	Name         string       `json:"name"`
	Type         int          `json:"type"`
	Amount       float64      `json:"amount"`
	Status       int          `json:"status"`
	PayAmount    float64      `json:"pay_amount"`
	RefundAmount float64      `json:"refund_amount"`
	PayStatus    int          `json:"pay_status"`
	PayTime      str.DateTime `json:"pay_time"`
	PayType      string       `json:"pay_type"`
	ExpiredTime  str.DateTime `json:"expired_time"`
	CompleteTime str.DateTime `json:"complete_time"`
	CreateTime   str.DateTime `json:"create_time"`
	UpdateTime   str.DateTime `json:"update_time"`
}

// Extra 订单扩展信息
type Extra struct {
	RpcService string `json:"rpc_service"` //支付成功后调用的rpc服务名
	RpcFunc    string `json:"rpc_func"`    //支付成功后调用的rpc服务方法名
	Param      string `json:"param"`       //其他参数
}
