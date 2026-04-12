package socket

// Message 消息格式
type Message struct {
	Id   string      `json:"id"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Time int64       `json:"time"`
}
