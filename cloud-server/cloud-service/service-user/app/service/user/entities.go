package user

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

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
	OuterOrderNo string `json:"outer_order_no"`
}
