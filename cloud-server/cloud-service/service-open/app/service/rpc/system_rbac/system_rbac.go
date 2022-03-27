package system_rbac

import (
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

// Check 判断rbac
func (*service) Check(userId int, roleId int, path string, method string) bool {
	client, err := rpc.New(group.CloudService, name.ServiceAdmin).Call("/rpc/rbac/check", &rpc.Param{
		"user_id": userId,
		"role_id": roleId,
		"path":    path,
		"method":  method,
	}).Response()
	if err != nil {
		return false
	}
	return client.Data.(bool)
}
