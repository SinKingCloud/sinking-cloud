package menu

// Type 菜单类型
type Type int

const (
	TypeDir     Type = iota // 目录
	TypePage                // 页面
	TypeControl             //空间
)

// Types 菜单类型数据
func Types() map[Type]string {
	return map[Type]string{
		TypeDir:     "目录",
		TypePage:    "页面",
		TypeControl: "控件",
	}
}
