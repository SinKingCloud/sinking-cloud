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

// AliyunSetting 阿里云设置
type AliyunSetting struct {
	AccessKeyID     string `json:"access_key_id"`     //AccessKeyId
	AccessKeySecret string `json:"access_key_secret"` //AccessKeySecret
	Bucket          string `json:"bucket"`            //Bucket名称
	Endpoint        string `json:"endpoint"`          //地域域名
	Path            string `json:"path"`              //保存路径
	Domain          string `json:"domain"`            //绑定域名
}
