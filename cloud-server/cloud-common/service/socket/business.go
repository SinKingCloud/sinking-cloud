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
func (*Service) SendToId(list []*SendToIdParam) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceSocket, "SendToId", list, &result)
	if err != nil {
		return false
	}
	return result
}

// SendToAll 发送全部消息
func (*Service) SendToAll(list []*SendToAllParam) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceSocket, "SendToAll", list, &result)
	if err != nil {
		return false
	}
	return result
}
