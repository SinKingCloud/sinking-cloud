package enum

import (
	"service-pay/app/enum/balance_log_type"
	"service-pay/app/enum/order_log_operator_type"
	"service-pay/app/enum/order_log_type"
	"service-pay/app/enum/order_pay_status"
	"service-pay/app/enum/order_pay_type"
	"service-pay/app/enum/order_status"
	"service-pay/app/enum/order_type"
)

// Data 枚举信息
var Data = map[string]interface{}{
	"order": map[string]interface{}{
		"status":     order_status.Map(),     //订单状态
		"type":       order_type.Map(),       //订单类型
		"pay_status": order_pay_status.Map(), //支付状态
		"pay_type":   order_pay_type.Map(),   //支付类型
	},
	"order_log": map[string]interface{}{
		"type":          order_log_type.Map(),          //订单日志类型
		"operator_type": order_log_operator_type.Map(), //订单日志操作人类型
	},
	"balance_log": map[string]interface{}{
		"type": balance_log_type.Map(), //余额日志操作类型
	},
}
