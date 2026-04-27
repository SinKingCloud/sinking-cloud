package channel_type

const (
	Alipay = "alipay"
	Wxpay  = "wxpay"
)

// Map 数据
func Map() map[string]string {
	return map[string]string{
		Alipay: "支付宝",
		Wxpay:  "微信",
	}
}
