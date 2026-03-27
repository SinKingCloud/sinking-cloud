package rpc

import (
	"encoding/json"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/user"
)

type userService struct {
}

// GetUserInfo 获取用户信息
func (*userService) GetUserInfo(params json.RawMessage) (interface{}, error) {
	var param *user.GetWebUserInfoParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return nil, err
	}
	u, err := service.User.GetUserInfo(param.UserId, param.Refresh)
	if err != nil {
		return nil, err
	}
	return u, nil
}
