package rpc

import (
	"encoding/json"
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/auth"
)

type authService struct {
}

// Login 账号登录
func (*authService) Login(params json.RawMessage) (interface{}, error) {
	var param *auth.WebUserLoginParam
	if err := json.Unmarshal(params, &param); err != nil {
		return "", err
	}
	token, err := service.Auth.Login(param.Account, param.Password, param.Code, param.UserAgent, param.Ip)
	if err != nil {
		return "", err
	}
	return token, nil
}
