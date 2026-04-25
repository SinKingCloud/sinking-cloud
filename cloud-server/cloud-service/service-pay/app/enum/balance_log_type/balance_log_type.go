package balance_log_type

const (
	Additional = iota
	Deduct
	Consume
	Recharge
	Freeze
	Defrost
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		Additional: "系统加款",
		Deduct:     "系统扣款",
		Consume:    "账户消费",
		Recharge:   "在线充值",
		Freeze:     "冻结余额",
		Defrost:    "解冻余额",
	}
}
