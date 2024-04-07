package log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/name"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
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

// CreateUser 创建网站用户日志
func (*Service) CreateUser(userId int, requestId string, requestIp string, types int, title string, content string) {
	go func(userId int, requestId string, requestIp string, types int, title string, content string) {
		_, _ = rpc.New(group.CloudService, name.ServiceUser).Call("/rpc/log/create", &rpc.Param{
			"user_id":    userId,
			"request_id": requestId,
			"request_ip": requestIp,
			"type":       types,
			"title":      title,
			"content":    content,
		}).Response()
	}(userId, requestId, requestIp, types, title, content)
}

// CreateAdmin 创建后台用户日志
func (*Service) CreateAdmin(userId int, requestId string, requestIp string, types int, title string, content string) {
	go func(userId int, requestId string, requestIp string, types int, title string, content string) {
		_, _ = rpc.New(group.CloudService, name.ServiceAdmin).Call("/rpc/log/create", &rpc.Param{
			"user_id":    userId,
			"request_id": requestId,
			"request_ip": requestIp,
			"type":       types,
			"title":      title,
			"content":    content,
		}).Response()
	}(userId, requestId, requestIp, types, title, content)
}
