package oauth

import (
	"errors"
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

func (*Service) AuthTicket(userId int, ticket string, redirectUrl string) error {
	client, err := rpc.New(group.CloudService, name.ServiceOpen).Call("/rpc/oauth/auth_ticket", &rpc.Param{
		"user_id":      userId,
		"ticket":       ticket,
		"redirect_url": redirectUrl,
	}).Response()
	if err != nil {
		return err
	}
	if client.Code == 200 {
		return nil
	}
	return errors.New(client.Message)
}
