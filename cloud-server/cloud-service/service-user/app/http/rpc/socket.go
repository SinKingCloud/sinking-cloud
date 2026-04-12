package rpc

import (
	"encoding/json"
	"service-user/app/service"
)

type socketService struct {
}

// GenToken 生成连接token
func (s *socketService) GenToken(params json.RawMessage) (interface{}, error) {
	var param string
	err := json.Unmarshal(params, &param)
	if err != nil {
		return "", err
	}
	token, err := service.Socket.GenToken(param)
	if err != nil {
		return "", err
	}
	return token, nil
}
