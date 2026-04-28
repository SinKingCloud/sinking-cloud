package balance_log_type

const (
	Additional = iota
	Deduct
	Consume
	Recharge
	Freeze
	Defrost
	Refunded
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		Additional: "后台加款",
		Deduct:     "后台扣款",
		Consume:    "账户消费",
		Recharge:   "在线充值",
		Freeze:     "冻结余额",
		Defrost:    "解冻余额",
		Refunded:   "订单退款",
	}
}
