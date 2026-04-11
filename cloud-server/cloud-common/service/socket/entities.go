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
