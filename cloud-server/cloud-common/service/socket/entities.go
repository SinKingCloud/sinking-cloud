package socket

// GenConnectionTokenParam 生成连接token
type GenConnectionTokenParam struct {
	Group string
	Id    string
}

// SendToIdParam 发送单播消息
type SendToIdParam struct {
	Group  string
	Id     string
	IsText bool
	Data   []byte
}

// SendToAllParam 发送全部消息
type SendToAllParam struct {
	Group  string
	IsText bool
	Data   []byte
}

// Message 前台 WebSocket 消息格式
type Message struct {
	Id   string      `json:"id"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Time int64       `json:"time"`
}

// SendToUserParam 给前台用户发送消息
type SendToUserParam struct {
	Message  *Message
	UserIds  []int64
	Sessions []string
}
