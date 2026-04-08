package role_status

const (
	Enable  = iota // 正常
	Disable        // 禁用
)

// Map 角色状态数据
func Map() map[int]string {
	return map[int]string{
		Enable:  "正常",
		Disable: "禁用",
	}
}
