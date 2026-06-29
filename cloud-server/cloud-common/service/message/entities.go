package message

// CreateMessageParam 创建消息参数
type CreateMessageParam struct {
	UserId  int64
	Title   string
	Content string
}

// CreateMessagesParam 批量创建消息参数
type CreateMessagesParam struct {
	UserIds []int64
	Title   string
	Content string
}
