package channel_type

const (
	Aliyun = iota // 正常
)

// Map 菜单状态数据
func Map() map[int]string {
	return map[int]string{
		Aliyun: "阿里云",
	}
}
