package rbac

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/constant"
)

func (s *service) Check(userId int, roleId int, path string, method string) bool {
	//判断权限信息
	if userId == constant.AdminId {
		return true
	}
	if r := s.CheckAccessByCache(roleId, path, method); r {
		return true
	} else {
		return false
	}
}
