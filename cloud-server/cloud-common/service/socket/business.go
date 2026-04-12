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

// BatchSendToId 分批发送单播消息
func (s *Service) BatchSendToId(list []*SendToIdParam, batch int) {
	if batch <= 0 {
		batch = 1000
	}
	if batch >= 100000 {
		batch = 100000
	}
	total := len(list)
	for i := 0; i < total; i += batch {
		end := i + batch
		if end > total {
			end = total
		}
		batchList := list[i:end]
		s.SendToId(batchList)
	}
}

// SendToAll 发送全部消息
func (*Service) SendToAll(list []*SendToAllParam) (result bool) {
	err := rpc.InvokeFunc(constant.ServiceSocket, "SendToAll", list, &result)
	if err != nil {
		return false
	}
	return result
}

// BatchSendToAll 分批发送全部消息
func (s *Service) BatchSendToAll(list []*SendToAllParam, batch int) {
	if batch <= 0 {
		batch = 1000
	}
	if batch >= 100000 {
		batch = 100000
	}
	total := len(list)
	for i := 0; i < total; i += batch {
		end := i + batch
		if end > total {
			end = total
		}
		batchList := list[i:end]
		s.SendToAll(batchList)
	}
}
