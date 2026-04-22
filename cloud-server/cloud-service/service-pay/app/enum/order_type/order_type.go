package order_type

const (
	Recharge = iota
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		Recharge: "在线充值",
	}
}
