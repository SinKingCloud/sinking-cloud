package rbac

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// RegisterPermissions 注册权限
func (*Service) RegisterPermissions(list []*RegisterPermissionsParam) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceAdmin, "RegisterPermissions", list, &result)
	if err != nil {
		return false
	}
	return result
}
