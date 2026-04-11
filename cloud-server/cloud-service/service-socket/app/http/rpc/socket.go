package rpc

import (
	"encoding/json"
	"service-socket/app/queue"
	"service-socket/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/socket"
)

type socketService struct {
}

// GenConnectionToken 生成连接token
func (*socketService) GenConnectionToken(params json.RawMessage) (interface{}, error) {
	var param *socket.GenConnectionTokenParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return "", err
	}
	token, err := service.Socket.GenConnectionToken(param.Group, param.Id)
	if err != nil {
		return "", err
	}
	return token, nil
}

// SendToId 发送单播消息
func (*socketService) SendToId(params json.RawMessage) (interface{}, error) {
	var param *socket.SendToIdParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return false, err
	}
	queue.Socket.SendToId(param.Group, param.Id, param.IsText, param.Data)
	return true, nil
}

// SendToAll 发送全部消息
func (*socketService) SendToAll(params json.RawMessage) (interface{}, error) {
	var param *socket.SendToAllParam
	err := json.Unmarshal(params, &param)
	if err != nil {
		return false, err
	}
	queue.Socket.SendToAll(param.Group, param.IsText, param.Data)
	return true, nil
}
