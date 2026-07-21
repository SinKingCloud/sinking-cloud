package mind_share_status

const (
	Enable  int = iota // 已开启
	Disable            // 已关闭
)

// Map 导图分享状态数据
func Map() map[int]string {
	return map[int]string{
		Enable:  "已开启",
		Disable: "已关闭",
	}
}
