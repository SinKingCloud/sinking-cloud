package session

// UserType 用户类型
type UserType int

const (
	UserTypeWeb   UserType = 0 // 网站用户
	UserTypeAdmin UserType = 1 // 系统用户
)

// UserTypes 用户类型数据
func UserTypes() map[UserType]string {
	return map[UserType]string{
		UserTypeWeb:   "网站用户",
		UserTypeAdmin: "系统用户",
	}
}
