package rpc

import (
	"encoding/json"
	"service-admin/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
)

type authService struct {
}

// Login 账号登录
func (*authService) Login(params json.RawMessage) (interface{}, error) {
	var param *auth.SystemUserLoginParam
	if err := json.Unmarshal(params, &param); err != nil {
		return "", err
	}
	token, err := service.Auth.Login(param.Account, param.Password, param.Code, param.UserAgent, param.Ip)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Logout 登录注销
func (*authService) Logout(params json.RawMessage) (interface{}, error) {
	var param string
	if err := json.Unmarshal(params, &param); err != nil {
		return false, err
	}
	err := service.Auth.Logout(param)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ClearExpireSessions 清理过期会话
func (*authService) ClearExpireSessions(_ json.RawMessage) (interface{}, error) {
	if err := service.Auth.ClearExpireSessions(); err != nil {
		return false, err
	}
	return true, nil
}
