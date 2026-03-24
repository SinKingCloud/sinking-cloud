package service

// RbacService service接口
type RbacService interface {
	CheckUserHasPermission(userId int64, path string) error
}
