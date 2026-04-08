package user_status

const (
	Enable  int = iota // 正常
	Disable            // 禁用
)

// Map 用户状态数据
func Map() map[int]string {
	return map[int]string{
		Enable:  "正常",
		Disable: "禁用",
	}
}
