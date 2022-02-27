package user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service"
)

// GetCurrentUserInfo  获取登录用户最新信息
func GetCurrentUserInfo(id int, refresh bool) *model.SystemUser {
	return service.GetSystemUserServiceIns().FindByIdWithCache(id, refresh)
}
