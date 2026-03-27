package service

// RbacService service接口
type RbacService interface {
	RegisterPermission(name string, code string, path string) error //注册权限
	CheckUserHasPermission(userId int64, path string) error         //判断用户是否含有权限
}
