package task_status

const (
	Running = iota // 正常
	Stop           // 禁用
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		Running: "运行",
		Stop:    "暂停",
	}
}
