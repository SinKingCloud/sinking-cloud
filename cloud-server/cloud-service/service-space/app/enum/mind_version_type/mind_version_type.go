package mind_version_type

const (
	Manual int = iota // 手动保存
	Auto              // 自动保存
)

// Map 导图版本类型数据
func Map() map[int]string {
	return map[int]string{
		Manual: "手动保存",
		Auto:   "自动保存",
	}
}
