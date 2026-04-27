package channel_plugin

const (
	Alipay = iota
	Epay
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		Alipay: "支付宝",
		Epay:   "易支付",
	}
}
