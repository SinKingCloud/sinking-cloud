package order_log_operator_type

const (
	System = iota
	User
	Admin
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		System: "系统",
		User:   "用户",
		Admin:  "管理员",
	}
}
