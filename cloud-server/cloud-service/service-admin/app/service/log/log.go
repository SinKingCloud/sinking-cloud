package log

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

// SystemLog 系统操作日志表
type SystemLog struct {
	Id         int          `json:"id"`
	UserId     int          `json:"user_id"`
	Account    string       `json:"account"`
	NickName   string       `json:"nick_name"`
	Avatar     string       `json:"avatar"`
	Email      string       `json:"email"`
	Phone      int64        `json:"phone"`
	RequestId  string       `json:"request_id"`
	RequestIp  string       `json:"request_ip"`
	Type       int          `json:"type"`
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}
