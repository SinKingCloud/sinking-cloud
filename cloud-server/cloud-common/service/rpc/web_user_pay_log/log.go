package web_user_pay_log

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

// Create 创建用户日志
func (*Service) Create(userId int, types int, logType int, money float64, title string, content string) {
	go func(userId int, types int, logType int, money float64, title string, content string) {
		_, _ = rpc.New(group.CloudService, name.ServiceUser).Call("/rpc/pay_log/create", &rpc.Param{
			"user_id":  userId,
			"type":     types,
			"log_type": logType,
			"money":    money,
			"title":    title,
			"content":  content,
		}).Response()
	}(userId, types, logType, money, title, content)
}
