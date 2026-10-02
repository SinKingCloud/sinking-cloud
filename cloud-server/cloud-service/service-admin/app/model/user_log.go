package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// UserLog 用户操作记录
type UserLog struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int64        `gorm:"column:user_id" json:"user_id"`
	RequestId  string       `gorm:"column:request_id" json:"request_id"`
	RequestIp  string       `gorm:"column:request_ip" json:"request_ip"`
	Location   string       `gorm:"column:location" json:"location"`
	Type       int          `gorm:"column:type" json:"type"`
	Title      string       `gorm:"column:title" json:"title"`
	Content    string       `gorm:"column:content" json:"content"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
}

// TableName 获取表名
func (*UserLog) TableName() string {
	return "cloud_user_logs"
}

// BeforeCreate 创建前
func (t *UserLog) BeforeCreate(_ *gorm.DB) error {
	if time.Time(t.CreateTime).IsZero() {
		t.CreateTime = str.DateTime(time.Now())
	}
	return nil
}
