package user

import "service-admin/app/repository/user"

// User 用户列表展示数据
type User struct {
	*user.User
	RoleIds []int64 `json:"role_ids"`
}
