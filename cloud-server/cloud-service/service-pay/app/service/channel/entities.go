package channel

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// Channel 存储通道
type Channel struct {
	Id         int64        `json:"id"`
	Type       int          `json:"type"`
	Name       string       `json:"name"`
	Config     interface{}  `json:"config"`
	Status     int          `json:"status"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}

// Alipay 阿里云设置
type Alipay struct {
	AppId      string `json:"app_id"`      //应用ID
	PrivateKey string `json:"private_key"` //私钥
	PublicKey  string `json:"public_key"`  //公钥
}

// Epay 易支付设置
type Epay struct {
	Pid string `json:"pid"` //商户ID
	Key string `json:"key"` //密钥
	Url string `json:"url"` //地址
}
