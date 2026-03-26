package service

// RbacService service接口
type RbacService interface {
	RegisterPermission(name string, code string, path string) error
	CheckUserHasPermission(userId int64, path string) error
}
