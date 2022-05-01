package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppOauthLog oauth记录列表
type AppOauthLog struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	OauthId    int          `gorm:"column:oauth_id" json:"oauth_id"`
	Ticket     string       `gorm:"column:ticket" json:"ticket"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (AppOauthLog) TableName() string {
	return "app_oauth_logs"
}

// BeforeCreate 创建前
func (t *AppOauthLog) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppOauthLog) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
