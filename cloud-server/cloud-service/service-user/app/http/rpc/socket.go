package rpc

import (
	"encoding/json"
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/socket"
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

// SendToUser 给前台用户发送 WebSocket 消息
func (s *socketService) SendToUser(params json.RawMessage) (interface{}, error) {
	var param socket.SendToUserParam
	if err := json.Unmarshal(params, &param); err != nil {
		return false, err
	}
	to := make([]interface{}, 0, len(param.UserIds)+len(param.Sessions))
	for _, id := range param.UserIds {
		to = append(to, id)
	}
	for _, session := range param.Sessions {
		to = append(to, session)
	}
	if len(param.ExcludeUserIds) > 0 || len(param.ExcludeSessions) > 0 {
		exclude := make([]interface{}, 0, len(param.ExcludeUserIds)+len(param.ExcludeSessions))
		for _, id := range param.ExcludeUserIds {
			exclude = append(exclude, id)
		}
		for _, session := range param.ExcludeSessions {
			exclude = append(exclude, session)
		}
		if err := service.Socket.SendMessageExcludeSessions(param.Message, exclude, to...); err != nil {
			return false, err
		}
		return true, nil
	}
	if err := service.Socket.SendMessage(param.Message, to...); err != nil {
		return false, err
	}
	return true, nil
}
