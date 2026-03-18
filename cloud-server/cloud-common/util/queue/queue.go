package queue

type Queue interface {
	// SendUnicast 发送单播消息（负载均衡，消费者随机接收）
	SendUnicast(topic string, body string) (mqMsg MqMsg, err error)
	// SendBroadcast 发送广播消息（所有订阅者都能收到）
	SendBroadcast(topic string, body string) (mqMsg MqMsg, err error)
	// ListenReceiveUnicast 监听接收单播消息
	ListenReceiveUnicast(topic string, receiveDo func(mqMsg MqMsg)) (err error)
	// ListenReceiveBroadcast 监听接收广播消息
	ListenReceiveBroadcast(topic string, receiveDo func(mqMsg MqMsg)) (err error)
	Close()
}

const (
	_ = iota
	SendUnicast
	ReceiveUnicast
	SendBroadcast
	ReceiveBroadcast
)

type MqMsg struct {
	RunType int    `json:"run_type"`
	Topic   string `json:"topic"`
	Body    string `json:"body"`
}
