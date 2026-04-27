package order_pay_type

const (
	Unset = iota
	Balance
	Alipay
	Wxpay
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		Unset:   "未支付",
		Balance: "账户余额",
		Alipay:  "支付宝",
		Wxpay:   "微信",
	}
}
