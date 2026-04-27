package channel

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// Channel 存储通道
type Channel struct {
	Id         int64        `json:"id"`
	Type       string       `json:"type"`
	Plugin     int          `json:"plugin"`
	Name       string       `json:"name"`
	Config     interface{}  `json:"config"`
	Status     int          `json:"status"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}

// Config 支付配置
type Config struct {
	Min     float64 `json:"min"`     //最小充值金额(0为不限制)
	Max     float64 `json:"max"`     //最大充值金额(0为不限制)
	Balance int64   `json:"balance"` //余额支付通道(0为关闭)
	Alipay  int64   `json:"alipay"`  //支付宝支付通道(0为关闭)
	Wxpay   int64   `json:"wxpay"`   //微信支付通道(0为关闭)
	Expire  int64   `json:"expire"`  //订单有效期(分钟)
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
