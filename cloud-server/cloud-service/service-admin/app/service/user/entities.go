package user

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// User 用户列表展示数据
type User struct {
	Id         int64        `json:"id"`
	Account    string       `json:"account"`
	Phone      string       `json:"phone"`
	Email      string       `json:"email"`
	Avatar     string       `json:"avatar"`
	NickName   string       `json:"nick_name"`
	Status     int          `json:"status"`
	RoleIds    []int64      `json:"role_ids"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}
