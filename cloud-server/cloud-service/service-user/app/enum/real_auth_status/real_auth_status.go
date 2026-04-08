package real_auth_status

const (
	UnAuth = iota // 未实名
	Authed        // 已实名
)

// Map 用户状态数据
func Map() map[int]string {
	return map[int]string{
		UnAuth: "未实名",
		Authed: "已实名",
	}
}
