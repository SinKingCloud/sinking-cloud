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

type Balance struct {
	Money         float64      `json:"money"`          //余额
	Frozen        float64      `json:"frozen"`         //冻结余额
	TotalRecharge float64      `json:"total_recharge"` //累计充值
	TotalConsume  float64      `json:"total_consume"`  //累计消费
	CreateTime    str.DateTime `json:"create_time"`    //创建时间
	UpdateTime    str.DateTime `json:"update_time"`    //更新时间
}
