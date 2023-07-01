package face

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/name"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"sync"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var (
	obj  *service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *service {
	once.Do(func() {
		obj = &service{}
	})
	return obj
}

// Check 判断验证码是否正确
func (*service) Check(id string) bool {
	client, err := rpc.New(group.CloudService, name.ServiceUser).Call("/rpc/face/check", &rpc.Param{
		"id": id,
	}).Response()
	if err != nil {
		return false
	}
	return client.Data.(bool)
}
