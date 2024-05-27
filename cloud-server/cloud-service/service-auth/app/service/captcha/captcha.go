package captcha

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/captcha"
	"sync"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// Check 判断验证码是否正确
func (*Service) Check(randStr string, ticket string) bool {
	return captcha.Tencent.Check(randStr, ticket)
}
