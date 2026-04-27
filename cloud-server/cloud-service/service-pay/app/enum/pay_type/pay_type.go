package pay_type

const (
	Balance = "balance"
	Alipay  = "alipay"
	Wxpay   = "wxpay"
)

// Map 数据
func Map() map[string]string {
	return map[string]string{
		Balance: "账户余额",
		Alipay:  "支付宝",
		Wxpay:   "微信",
	}
}
