package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Log 计划任务表
type Log struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	TaskId     int64        `gorm:"column:task_id" json:"task_id"`
	Content    string       `gorm:"column:content" json:"content"`
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
