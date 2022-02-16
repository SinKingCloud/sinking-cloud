package person

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

type ControllerPerson struct {
}

// Info 个人信息
type Info struct {
	Account   string       `json:"account"`
	NickName  string       `json:"nick_name"`
	Avatar    string       `json:"avatar"`
	Email     string       `json:"email"`
	Phone     int64        `json:"phone"`
	LoginIp   string       `json:"login_ip"`
	LoginTime str.DateTime `json:"login_time"`
}
