package menu_type

const (
	Dir     = iota // 目录
	Page           // 页面
	Control        // 控件
)

// Types 菜单类型数据
func Types() map[int]string {
	return map[int]string{
		Dir:     "目录",
		Page:    "页面",
		Control: "控件",
	}
}
