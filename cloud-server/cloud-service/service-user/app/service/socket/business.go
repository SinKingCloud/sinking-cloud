package socket

import (
	"encoding/json"
	"errors"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/socket"
)

// GenToken 生成连接token
func (s *service) GenToken(session string) (string, error) {
	info := s.authService.GetSessionInfo(session)
	if info == nil {
		return "", errors.New("获取会话信息失败")
	}
	token, err := s.socketService.GenConnectionToken(constant.ServiceUser, info.Token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// SendMessage 发送消息
func (s *service) SendMessage(message *socket.Message, to ...interface{}) error {
	if message == nil {
		return errors.New("消息不能为空")
	}
	bytes, err := json.Marshal(message)
	if err != nil {
		return errors.New("处理消息失败")
	}
	if len(to) <= 0 {
		s.socketService.BatchSendToAll([]*socket.SendToAllParam{{
			Group:  constant.ServiceUser,
			IsText: true,
			Data:   bytes,
		}}, 5000)
		return nil
	}
	sessions, err := s.parseSessions(to...)
	if err != nil {
		return err
	}
	return s.sendToSessions(bytes, sessions)
}

// sendToSessions 向指定会话发送消息，并对会话去重。
func (s *service) sendToSessions(bytes []byte, sessions []string) error {
	var list []*socket.SendToIdParam
	sessionMap := make(map[string]struct{}, len(sessions))
	for _, session := range sessions {
		if session == "" {
			continue
		}
		if _, ok := sessionMap[session]; ok {
			continue
		}
		sessionMap[session] = struct{}{}
		list = append(list, &socket.SendToIdParam{
			Group:  constant.ServiceUser,
			Id:     session,
			IsText: true,
			Data:   bytes,
		})
	}
	s.socketService.BatchSendToId(list, 5000)
	return nil
}

// parseSessions 将用户ID或会话转换为会话列表，支持 int/int64/[]int/[]int64/string/[]string/[]interface{}。
func (s *service) parseSessions(to ...interface{}) ([]string, error) {
	var ids []int64
	var sessions []string
	for _, item := range to {
		switch v := item.(type) {
		case nil:
			continue
		case int:
			ids = append(ids, int64(v))
		case []int:
			for _, id := range v {
				ids = append(ids, int64(id))
			}
		case int64:
			ids = append(ids, v)
		case []int64:
			ids = append(ids, v...)
		case string:
			sessions = append(sessions, v)
		case []string:
			sessions = append(sessions, v...)
		case []interface{}:
			for _, value := range v {
				itemSessions, err := s.parseSessions(value)
				if err != nil {
					return nil, err
				}
				sessions = append(sessions, itemSessions...)
			}
		default:
			return nil, errors.New("接收对象类型不支持")
		}
	}
	if len(ids) > 0 {
		tokenList, err := s.repositorySession.SelectUnExpireTokensByUserIds(ids)
		if err != nil {
			return nil, errors.New("获取用户会话失败")
		}
		sessions = append(sessions, tokenList...)
	}
	return sessions, nil
}
