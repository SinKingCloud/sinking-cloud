package order_status

const (
	UnPay = iota
	Payed
	Closed
	Refunded
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		UnPay:    "待支付",
		Payed:    "已支付",
		Closed:   "已关闭",
		Refunded: "已退款",
	}
}
