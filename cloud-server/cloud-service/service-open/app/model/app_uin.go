package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppUin QQ记录
type AppUin struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	AppId      int          `gorm:"column:app_id" json:"app_id"`
	Server     string       `gorm:"column:server" json:"server"`
	Uin        int64        `gorm:"column:uin" json:"uin"`
	Password   string       `gorm:"column:password" json:"password"`
	UinStatus  int          `gorm:"column:uin_status" json:"uin_status"`
	ExpireDay  int          `gorm:"column:expire_day" json:"expire_day"`
	PlayerDay  int          `gorm:"column:player_day" json:"player_day"`
	Setting    string       `gorm:"column:setting" json:"setting"`
	Timing     string       `gorm:"column:timing" json:"timing"`
	Status     int          `gorm:"column:status" json:"status"`
	RunTime    str.DateTime `gorm:"column:run_time" json:"run_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (AppUin) TableName() string {
	return "app_uins"
}

// BeforeCreate 创建前
func (t *AppUin) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppUin) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
