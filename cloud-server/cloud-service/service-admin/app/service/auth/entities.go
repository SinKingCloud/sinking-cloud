package auth

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// Session 会话表
type Session struct {
	Token      string       `gorm:"column:token" json:"token"`
	Device     string       `gorm:"column:device" json:"device"`
	Ip         string       `gorm:"column:ip" json:"ip"`
	Location   string       `gorm:"column:location" json:"location"`
	ExpireTime str.DateTime `gorm:"column:expire_time" json:"expire_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}
