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
func (s *service) SendMessage(message *Message, to ...interface{}) error {
	if message == nil {
		return errors.New("消息不能为空")
	}
	bytes, err := json.Marshal(message)
	if err != nil {
		return errors.New("处理消息失败")
	}
	if len(to) <= 0 {
		if s.socketService.SendToAll([]*socket.SendToAllParam{{
			Group:  constant.ServiceUser,
			IsText: true,
			Data:   bytes,
		}}) {
			return nil
		}
		return errors.New("发送消息失败")
	}
	var ids []int64
	var sessions []string
	for _, item := range to {
		switch v := item.(type) {
		case int64:
			ids = append(ids, v)
		case []int64:
			ids = append(ids, v...)
		case string:
			sessions = append(sessions, v)
		case []string:
			sessions = append(sessions, v...)
		default:
			return errors.New("接收对象类型不支持")
		}
	}
	if len(ids) > 0 {
		tokenList, err2 := s.repositorySession.SelectUnExpireTokensByUserIds(ids)
		if err2 != nil {
			return errors.New("获取用户会话失败")
		}
		sessions = append(sessions, tokenList...)
	}
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
	const batchSize = 5000
	total := len(list)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batchList := list[i:end]
		if !s.socketService.SendToId(batchList) {
			return errors.New("发送消息失败")
		}
	}
	return nil
}
