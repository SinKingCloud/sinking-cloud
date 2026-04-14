package user

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// SelectUser 用户查询
type SelectUser struct {
	Account         string
	Phone           string
	Email           string
	NickName        string
	Status          string
	CreateTimeStart string
	CreateTimeEnd   string
	UpdateTimeStart string
	UpdateTimeEnd   string
}

// UpdateUser 用户更新
type UpdateUser struct {
	Account  interface{}
	Phone    interface{}
	Email    interface{}
	Password interface{}
	Avatar   interface{}
	NickName interface{}
	Status   interface{}
}

// User 用户列表展示数据
type User struct {
	Id         int64        `json:"id"`
	Account    string       `json:"account"`
	Phone      string       `json:"phone"`
	Email      string       `json:"email"`
	Avatar     string       `json:"avatar"`
	NickName   string       `json:"nick_name"`
	Status     int          `json:"status"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}
