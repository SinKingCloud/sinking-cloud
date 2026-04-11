package socket

import "errors"

// SendToLocalId 给指定group指定id发送消息
func (s *service) SendToLocalId(group string, id string, messageType int, payload []byte) error {
	connection, err := s.getConnection(group, id)
	if err != nil {
		return err
	}
	err = connection.Send(messageType, payload)
	if err != nil {
		return errors.New("发送失败")
	}
	return nil
}

// SendToLocalAll 给指定group所有id发送消息
func (s *service) SendToLocalAll(group string, messageType int, payload []byte) (queued int, dropped int, closed int, err error) {
	registry := s.getConnectionPool(group)
	result, err := registry.Broadcast(messageType, payload)
	if err != nil {
		return 0, 0, 0, err
	}
	return result.Queued, result.Dropped, result.Closed, nil
}
