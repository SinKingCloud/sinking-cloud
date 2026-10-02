package user_log

// SelectUserLog 用户操作日志查询
type SelectUserLog struct {
	UserId          *int64
	Keyword         *string
	RequestId       *string
	RequestIp       *string
	Location        *string
	Type            *int
	Title           *string
	Content         *string
	CreateTimeStart *string
	CreateTimeEnd   *string
}
