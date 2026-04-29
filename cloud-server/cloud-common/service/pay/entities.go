package pay

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

type NotifyParam struct {
	UserId     int64        `json:"user_id"`      //用户ID
	TradeNo    string       `json:"trade_no"`     //内部订单号
	Extra      string       `json:"extra"`        //扩展信息
	Amount     float64      `json:"amount"`       //订单金额
	OutTradeNo string       `json:"out_trade_no"` //外部订单号
	PayAmount  float64      `json:"pay_amount"`   //支付金额
	PayType    int          `json:"pay_type"`     //支付方式
	PayTime    str.DateTime `json:"pay_time"`     //支付时间
}
