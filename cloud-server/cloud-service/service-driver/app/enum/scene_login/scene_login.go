package scene_login

const (
	Yes = iota
	No
)

// Map 菜单状态数据
func Map() map[int]string {
	return map[int]string{
		Yes: "需要登录",
		No:  "无需登录",
	}
}
