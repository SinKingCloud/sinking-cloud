package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

type Job struct {
	Id          int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Code        int          `gorm:"column:code" json:"code"`
	Status      int          `gorm:"column:status" json:"status"`
	CreateTime  str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime  str.DateTime `gorm:"column:update_time" json:"update_time"`
	Subject     int          `gorm:"column:subject" json:"subject"`
	LastIndex   int          `gorm:"column:last_index" json:"last_index"`
	LastRunTime str.DateTime `gorm:"column:last_run_time" json:"last_run_time"`
}

// TableName 获取表名
func (Job) TableName() string {
	return "spider_jobs"
}

// BeforeCreate 创建前
func (t *Job) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Job) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
