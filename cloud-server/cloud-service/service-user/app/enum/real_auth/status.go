package real_auth

// Statu 用户状态
type Statu int

const (
	UnAuth Statu = iota // 未实名
	Authed              // 已实名
)

// Status 用户状态数据
func Status() map[Statu]string {
	return map[Statu]string{
		UnAuth: "未实名",
		Authed: "已实名",
	}
}
