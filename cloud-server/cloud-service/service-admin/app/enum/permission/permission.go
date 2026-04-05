package permission

import "service-admin/app/service"

// CodeMap 权限标识
func CodeMap() map[string]string {
	data, err := service.Rbac.GetPermissionEnum(false)
	if err != nil {
		return map[string]string{}
	}
	return data
}
