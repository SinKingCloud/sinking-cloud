package user_log_type

const (
	EventLogin = iota
	EventShow
	EventDelete
	EventEdit
	EventCreate
)

// Map 数据
func Map() map[int]string {
	return map[int]string{
		EventLogin:  "登陆事件",
		EventShow:   "查看事件",
		EventDelete: "删除事件",
		EventEdit:   "修改事件",
		EventCreate: "创建事件",
	}
}
