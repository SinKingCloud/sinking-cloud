package rpc

import (
	"encoding/json"
	"errors"
	"service-user/app/service"

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
	u, err := service.User.GetInfo(param.UserId, param.Refresh)
	if err != nil {
		return nil, err
	}
	if u == nil || u.Id <= 0 {
		return nil, errors.New("用户不存在")
	}
	return u, nil
}
