package rpc

import (
	"encoding/json"
	"service-user/app/service"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/message"
)

type messageService struct {
}

// Create 创建消息
func (*messageService) Create(params json.RawMessage) (interface{}, error) {
	var param message.CreateMessageParam
	if err := json.Unmarshal(params, &param); err != nil {
		return false, err
	}
	if err := service.Message.Create(param.UserId, param.Title, param.Content); err != nil {
		return false, err
	}
	return true, nil
}

// Creates 批量创建消息
func (*messageService) Creates(params json.RawMessage) (interface{}, error) {
	var param message.CreateMessagesParam
	if err := json.Unmarshal(params, &param); err != nil {
		return false, err
	}
	if err := service.Message.Creates(param.UserIds, param.Title, param.Content); err != nil {
		return false, err
	}
	return true, nil
}
