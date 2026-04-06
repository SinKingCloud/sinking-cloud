package role_name

import "service-admin/app/service"

// NameMap 角色名称
func NameMap() map[int64]string {
	data, err := service.Rbac.GetRoleNameEnum(false)
	if err != nil {
		return map[int64]string{}
	}
	return data
}
