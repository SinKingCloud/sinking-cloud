package order_pay_type

const (
	Balance = iota
	Alipay
	Wxpay
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		Balance: "账户余额",
		Alipay:  "支付宝",
		Wxpay:   "微信",
	}
}
