package user

import (
	repositoryUser "service-user/app/repository/user"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// User 用户信息
type User struct {
	*repositoryUser.User
	Money         float64 `json:"money"`          //余额
	Frozen        float64 `json:"frozen"`         //冻结余额
	TotalRecharge float64 `json:"total_recharge"` //累计充值
	TotalConsume  float64 `json:"total_consume"`  //累计消费
}

// RealAuthFiles 实名认证资料信息
type RealAuthFiles struct {
	Bucket    string   `json:"bucket"`
	PicPaths  []string `json:"pic_paths"`
	VideoPath string   `json:"video_path"`
}

// RealAuthInfo 实名认证信息
type RealAuthInfo struct {
	UserId      int64          `json:"user_id"`
	Type        int            `json:"type"`
	Status      int            `json:"status"`
	CompanyName string         `json:"company_name"`
	CompanyNo   string         `json:"company_no"`
	CardName    string         `json:"card_name"`
	CardNo      string         `json:"card_no"`
	Files       *RealAuthFiles `json:"files"`
	AuthTime    str.DateTime   `json:"auth_time"`
}

// realAuthProcess 实名认证流程缓存
type realAuthProcess struct {
	UserId       int64  `json:"user_id"`
	Type         int    `json:"type"`
	CompanyName  string `json:"company_name"`
	CompanyNo    string `json:"company_no"`
	CardName     string `json:"card_name"`
	CardNo       string `json:"card_no"`
	CertifyId    string `json:"certify_id"`
	CertifyUrl   string `json:"certify_url"`
	SceneId      int64  `json:"scene_id"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	OuterOrderNo string `json:"outer_order_no"`
}
