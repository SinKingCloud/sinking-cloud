package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppOauthUser oauth用户列表
type AppOauthUser struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	OauthId    int          `gorm:"column:oauth_id" json:"oauth_id"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	Status     int          `gorm:"column:status" json:"status"`
	LoginTime  str.DateTime `gorm:"column:login_time" json:"login_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (AppOauthUser) TableName() string {
	return "app_oauth_users"
}

// BeforeCreate 创建前
func (t *AppOauthUser) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppOauthUser) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
