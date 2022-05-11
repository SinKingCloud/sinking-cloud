package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppUinInfo QQ记录
type AppUinInfo struct {
	Id          int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Uin         string       `gorm:"column:uin" json:"uin"`
	NickName    string       `gorm:"column:nick_name" json:"nick_name"`
	Level       int          `gorm:"column:level" json:"level"`
	ActivateDay string       `gorm:"column:activate_day" json:"activate_day"`
	TodayDay    string       `gorm:"column:today_day" json:"today_day"`
	SpeedInfo   string       `gorm:"column:speed_info" json:"speed_info"`
	Skey        string       `gorm:"column:skey" json:"skey"`
	PSkey       string       `gorm:"column:p_skey" json:"p_skey"`
	SuperToken  string       `gorm:"column:super_token" json:"super_token"`
	CreateTime  str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime  str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (AppUinInfo) TableName() string {
	return "app_uin_infos"
}

// BeforeCreate 创建前
func (t *AppUinInfo) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppUinInfo) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
