package session

// SelectSession 会话查询
type SelectSession struct {
	UserId          *int64
	Keyword         *string
	Device          *string
	Ip              *string
	Location        *string
	CreateTimeStart *string
	CreateTimeEnd   *string
	UpdateTimeStart *string
	UpdateTimeEnd   *string
	ExpireTimeStart *string
	ExpireTimeEnd   *string
}
