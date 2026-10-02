package rpc

import (
	"encoding/json"
	"errors"
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
	u, err := service.User.GetInfo(param.UserId, param.Refresh)
	if err != nil {
		return nil, err
	}
	if u == nil || u.Id <= 0 {
		return nil, errors.New("用户不存在")
	}
	return u, nil
}

// CreateUserLog 异步创建用户操作日志
func (*userService) CreateUserLog(params json.RawMessage) (interface{}, error) {
	var param user.CreateSystemUserLogParam
	if err := json.Unmarshal(params, &param); err != nil {
		return false, err
	}
	if err := service.User.CreateLog(param.Id, param.UserId, param.RequestId, param.RequestIp, param.Type, param.Title, param.Content); err != nil {
		return false, err
	}
	return true, nil
}
