package person

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

type ControllerPerson struct {
}

// Info 个人信息
type Info struct {
	Id        int          `json:"id"`
	Account   string       `json:"account"`
	Money     float64      `json:"money"`
	RealAuth  int          `json:"real_auth"`
	NickName  string       `json:"nick_name"`
	Avatar    string       `json:"avatar"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	LoginIp   string       `json:"login_ip"`
	LoginTime str.DateTime `json:"login_time"`
}
