package socket

import (
	"errors"
	"fmt"

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

// SendToUser 给前台用户发送消息；to 支持 int/int64/[]int/[]int64/string/[]string，不传则广播给全部前台用户连接。
func (s *Service) SendToUser(message *Message, to ...interface{}) error {
	if message == nil {
		return errors.New("消息不能为空")
	}
	param := &SendToUserParam{Message: message}
	var fillParam func(...interface{}) error
	fillParam = func(items ...interface{}) error {
		for _, item := range items {
			switch v := item.(type) {
			case nil:
				continue
			case int:
				param.UserIds = append(param.UserIds, int64(v))
			case []int:
				for _, id := range v {
					param.UserIds = append(param.UserIds, int64(id))
				}
			case int64:
				param.UserIds = append(param.UserIds, v)
			case []int64:
				param.UserIds = append(param.UserIds, v...)
			case string:
				param.Sessions = append(param.Sessions, v)
			case []string:
				param.Sessions = append(param.Sessions, v...)
			case []interface{}:
				if err := fillParam(v...); err != nil {
					return err
				}
			default:
				return fmt.Errorf("接收对象类型不支持: %T", item)
			}
		}
		return nil
	}
	if err := fillParam(to...); err != nil {
		return err
	}
	var result bool
	err := rpc.InvokeFunc(constant.ServiceUser, "SendToUser", param, &result)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("发送消息失败")
	}
	return nil
}
