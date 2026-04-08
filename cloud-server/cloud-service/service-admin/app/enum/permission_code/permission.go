package permission_code

import "service-admin/app/service"

// Map 权限标识
func Map() map[string]string {
	data, err := service.Rbac.GetPermissionEnum(false)
	if err != nil {
		return map[string]string{}
	}
	return data
}
