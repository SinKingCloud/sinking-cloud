package rpc

import (
	"encoding/json"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rbac"
)

type rbacService struct {
}

// RegisterPermissions 注册权限
func (*rbacService) RegisterPermissions(params json.RawMessage) (interface{}, error) {
	var list []*rbac.RegisterPermissionsParam
	err := json.Unmarshal(params, &list)
	if err != nil {
		return false, err
	}
	if list != nil && len(list) > 0 {
		for _, v := range list {
			service.Rbac.RegisterPermissionAsync(v.Name, v.Code, v.Path)
		}
	}
	return true, nil
}
