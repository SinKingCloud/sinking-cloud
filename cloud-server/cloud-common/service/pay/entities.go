package pay

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

type NotifyParam struct {
	UserId    int64        `json:"user_id"`    //用户ID
	OrderId   int64        `json:"order_id"`   //订单ID
	Extra     string       `json:"extra"`      //扩展信息
	Amount    float64      `json:"amount"`     //订单金额
	TradeNo   string       `json:"trade_no"`   //第三方支付订单号
	PayAmount float64      `json:"pay_amount"` //支付金额
	PayType   string       `json:"pay_type"`   //支付方式
	PayTime   str.DateTime `json:"pay_time"`   //支付时间
}

type CreateOrderParam struct {
	UserId             int64   `json:"user_id"`              //用户ID
	Name               string  `json:"name"`                 //订单名称
	Describe           string  `json:"describe"`             //订单描述
	Amount             float64 `json:"amount"`               //订单金额
	PayType            string  `json:"pay_type"`             //支付类型
	Param              string  `json:"param"`                //业务扩展参数
	CallbackRpcService string  `json:"callback_rpc_service"` //支付回调rpc服务
	CallbackRpcFunc    string  `json:"callback_rpc_func"`    //支付类型rpc方法
}
