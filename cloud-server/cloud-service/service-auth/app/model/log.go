package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Log 日志表
type Log struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int64        `gorm:"column:user_id" json:"user_id"`
	UserType   int          `gorm:"column:user_type" json:"user_type"`
	Type       int          `gorm:"column:type" json:"type"`
	Describe   string       `gorm:"column:describe" json:"describe"`
	Ip         string       `gorm:"column:ip" json:"ip"`
	Location   string       `gorm:"column:location" json:"location"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
}

// TableName 获取表名
func (*Log) TableName() string {
	return "cloud_logs"
}

// BeforeCreate 创建前
func (t *Log) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}
