package menu

// Statu 菜单状态
type Statu int

const (
	StatusEnable  Statu = iota // 正常
	StatusDisable              // 禁用
)

// Status 用户状态数据
func Status() map[Statu]string {
	return map[Statu]string{
		StatusEnable:  "正常",
		StatusDisable: "禁用",
	}
}
