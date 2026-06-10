package order_status

const (
	UnPay = iota
	Delivering
	Credited
	Refunded
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		UnPay:      "待支付",
		Delivering: "发货中",
		Credited:   "已到账",
		Refunded:   "已退款",
	}
}
