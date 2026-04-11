package socket

// SendToId 发送单播消息
func (s *Queue) SendToId(group string, id string, isText bool, data []byte) {
	s.searchUnicast.SendUnicastTask(&Message{
		Group:  group,
		Id:     id,
		IsText: isText,
		Data:   data,
	})
}

// SendToAll 发送全部消息
func (s *Queue) SendToAll(group string, isText bool, data []byte) {
	s.sendBroadcast.SendBroadcastTask(&Message{
		Group:  group,
		IsText: isText,
		Data:   data,
	})
}
