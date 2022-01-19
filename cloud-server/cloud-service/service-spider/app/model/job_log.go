package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// JobLog 爬虫任务表
type JobLog struct {
	Id          int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	CityCode    string       `gorm:"column:city_code" json:"city_code"`
	SubjectCode string       `gorm:"column:subject_code" json:"subject_code"`
	Index       string       `gorm:"column:index" json:"index"`
	RegCode     string       `gorm:"column:reg_code" json:"reg_code"`
	TycId       string       `gorm:"column:tyc_id" json:"tyc_id"`
	CreditCode  string       `gorm:"column:credit_code" json:"credit_code"`
	CompanyName string       `gorm:"column:company_name" json:"company_name"`
	Content     string       `gorm:"column:content" json:"content"`
	Status      int          `gorm:"column:status" json:"status"`
	CreateTime  str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime  str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (JobLog) TableName() string {
	return "spider_job_logs"
}

// BeforeCreate 创建前
func (t *JobLog) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *JobLog) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
