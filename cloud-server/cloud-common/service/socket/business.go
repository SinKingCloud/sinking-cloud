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
	if err := s.fillSendToUserParam(param, false, to...); err != nil {
		return err
	}
	return s.sendToUser(param)
}

// SendToUserExcludeSessions 给指定前台用户或会话发送消息，并排除指定会话或用户；exclude 支持 int/int64/[]int/[]int64/string/[]string。
func (s *Service) SendToUserExcludeSessions(message *Message, exclude interface{}, to ...interface{}) error {
	if message == nil {
		return errors.New("消息不能为空")
	}
	param := &SendToUserParam{Message: message}
	if err := s.fillSendToUserParam(param, false, to...); err != nil {
		return err
	}
	if len(param.UserIds) == 0 && len(param.Sessions) == 0 {
		return errors.New("排除发送必须指定接收对象")
	}
	if err := s.fillSendToUserParam(param, true, exclude); err != nil {
		return err
	}
	return s.sendToUser(param)
}

func (s *Service) fillSendToUserParam(param *SendToUserParam, exclude bool, to ...interface{}) error {
	for _, item := range to {
		switch v := item.(type) {
		case nil:
			continue
		case int:
			if exclude {
				param.ExcludeUserIds = append(param.ExcludeUserIds, int64(v))
			} else {
				param.UserIds = append(param.UserIds, int64(v))
			}
		case []int:
			for _, id := range v {
				if exclude {
					param.ExcludeUserIds = append(param.ExcludeUserIds, int64(id))
				} else {
					param.UserIds = append(param.UserIds, int64(id))
				}
			}
		case int64:
			if exclude {
				param.ExcludeUserIds = append(param.ExcludeUserIds, v)
			} else {
				param.UserIds = append(param.UserIds, v)
			}
		case []int64:
			if exclude {
				param.ExcludeUserIds = append(param.ExcludeUserIds, v...)
			} else {
				param.UserIds = append(param.UserIds, v...)
			}
		case string:
			if exclude {
				param.ExcludeSessions = append(param.ExcludeSessions, v)
			} else {
				param.Sessions = append(param.Sessions, v)
			}
		case []string:
			if exclude {
				param.ExcludeSessions = append(param.ExcludeSessions, v...)
			} else {
				param.Sessions = append(param.Sessions, v...)
			}
		case []interface{}:
			for _, value := range v {
				if err := s.fillSendToUserParam(param, exclude, value); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("接收对象类型不支持: %T", item)
		}
	}
	return nil
}

func (s *Service) sendToUser(param *SendToUserParam) error {
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
