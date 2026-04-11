package socket

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

// GenConnectionToken 生成授权token
func (*Service) GenConnectionToken(group string, id string) (token string, err error) {
	err = rpc.InvokeFunc(constant.ServiceSocket, "GenConnectionToken", &GenConnectionTokenParam{
		Group: group,
		Id:    id,
	}, &token)
	if err != nil {
		return
	}
	return token, nil
}

// SendToId 发送单播消息
func (*Service) SendToId(group string, id string, isText bool, data []byte) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceSocket, "SendToId", &SendToIdParam{
		Group:  group,
		Id:     id,
		IsText: isText,
		Data:   data,
	}, &result)
	if err != nil {
		return false
	}
	return result
}

// SendToAll 发送全部消息
func (*Service) SendToAll(group string, isText bool, data []byte) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceSocket, "SendToAll", &SendToAllParam{
		Group:  group,
		IsText: isText,
		Data:   data,
	}, &result)
	if err != nil {
		return false
	}
	return result
}
