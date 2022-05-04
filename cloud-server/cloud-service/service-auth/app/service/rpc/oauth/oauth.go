package oauth

import (
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/name"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

func (*service) AuthTicket(userId int, ticket string, redirectUrl string) error {
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
