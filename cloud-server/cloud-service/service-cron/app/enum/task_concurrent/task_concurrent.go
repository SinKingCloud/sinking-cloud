package task_concurrent

const (
	Enable  = iota // 正常
	Disable        // 禁用
)

// Map 映射数据
func Map() map[int]string {
	return map[int]string{
		Enable:  "允许",
		Disable: "禁止",
	}
}
