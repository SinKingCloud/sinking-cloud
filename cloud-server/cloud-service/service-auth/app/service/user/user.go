package user

// Service 单例对象
type Service struct {
}

// obj 单例对象
var obj *Service

// GetIns 获取单例
func GetIns() *Service {
	if obj == nil {
		obj = &Service{}
	}
	return obj
}
