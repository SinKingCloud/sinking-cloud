package order_log_type

const (
	Create = iota
	Pay
	Refunded
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		Create:   "创建订单",
		Pay:      "支付订单",
		Refunded: "订单退款",
	}
}
