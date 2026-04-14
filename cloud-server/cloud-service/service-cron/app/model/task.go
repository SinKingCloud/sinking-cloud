package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Task 计划任务表
type Task struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Type       int          `gorm:"column:type" json:"type"`
	Name       string       `gorm:"column:name" json:"name"`
	Spec       string       `gorm:"column:spec" json:"spec"`
	Script     string       `gorm:"column:script" json:"script"`
	Concurrent int          `gorm:"column:concurrent" json:"concurrent"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Task) TableName() string {
	return "cloud_tasks"
}

// BeforeCreate 创建前
func (t *Task) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Task) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
