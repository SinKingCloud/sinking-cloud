package channel_status

const (
	Enable  = iota // 正常
	Disable        // 禁用
)

// Map 映射数据
func Map() map[int]string {
	return map[int]string{
		Enable:  "正常",
		Disable: "禁用",
	}
}
