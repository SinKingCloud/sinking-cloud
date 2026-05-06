package real_auth_type

const (
	UnAuth = iota
	Person
	Company
)

// Map 用户状态数据
func Map() map[int]string {
	return map[int]string{
		UnAuth:  "未实名",
		Person:  "个人实名",
		Company: "企业实名",
	}
}
