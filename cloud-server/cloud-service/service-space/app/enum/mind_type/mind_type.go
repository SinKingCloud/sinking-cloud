package mind_type

const (
	All        int = iota // 全部
	Product               // 产品构想
	Technology            // 技术方案
	Knowledge             // 知识整理
	Personal              // 个人计划
	Other                 // 其他
)

// Map 导图预设分类数据
func Map() map[int]string {
	return map[int]string{
		All:        "全部",
		Product:    "产品构想",
		Technology: "技术方案",
		Knowledge:  "知识整理",
		Personal:   "个人计划",
		Other:      "其他",
	}
}
