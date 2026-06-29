package message_status

const (
	Unread int = iota // 未读
	Read              // 已读
)

// Map 消息状态数据
func Map() map[int]string {
	return map[int]string{
		Unread: "未读",
		Read:   "已读",
	}
}
