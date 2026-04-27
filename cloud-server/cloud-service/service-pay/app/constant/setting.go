package constant

const (
	PayGroup   = "pay"                 //支付配置
	PayMin     = PayGroup + ".min"     //最小充值金额
	PayMax     = PayGroup + ".max"     //最大充值金额
	PayBalance = PayGroup + ".balance" //是否开启余额支付(0:关闭,1:开启)
	PayAlipay  = PayGroup + ".alipay"  //是否开启支付宝支付(0:关闭,其他:通道ID)
	PayWxpay   = PayGroup + ".wxpay"   //是否开启微信支付(0:关闭,其他:通道ID)
	PayExpire  = PayGroup + ".expire"  //订单有效期(分钟)
)
