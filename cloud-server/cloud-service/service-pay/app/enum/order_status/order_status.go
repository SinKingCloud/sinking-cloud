package order_status

const (
	UnPay = iota
	Payed
	Refunded
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		UnPay:    "待支付",
		Payed:    "已支付",
		Refunded: "已退款",
	}
}
