package task_type

const (
	Script = iota //脚本
	Rpc           //rpc调用
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		Script: "系统脚本",
		Rpc:    "调用方法",
	}
}
