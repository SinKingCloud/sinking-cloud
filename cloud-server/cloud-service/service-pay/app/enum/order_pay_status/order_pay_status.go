package order_pay_status

const (
	UnPay = iota
	Payed
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		UnPay: "未支付",
		Payed: "已支付",
	}
}
