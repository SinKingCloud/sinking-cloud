package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// SystemLogRequest 请求日志表
type SystemLogRequest struct {
	Id             string       `gorm:"column:id;PRIMARY_KEY" json:"id"`
	RequestId      string       `gorm:"column:request_id" json:"request_id"`
	RequestIp      string       `gorm:"column:request_ip" json:"request_ip"`
	RequestUrl     string       `gorm:"column:request_url" json:"request_url"`
	RequestMethod  string       `gorm:"column:request_method" json:"request_method"`
	RequestHeader  string       `gorm:"column:request_header" json:"request_header"`
	RequestBody    string       `gorm:"column:request_body" json:"request_body"`
	RequestTime    str.DateTime `gorm:"column:request_time" json:"request_time"`
	ResponseCode   int          `gorm:"column:response_code" json:"response_code"`
	ResponseHeader string       `gorm:"column:response_header" json:"response_header"`
	ResponseBody   string       `gorm:"column:response_body" json:"response_body"`
	ResponseTime   int          `gorm:"column:response_time" json:"response_time"`
	GroupName      string       `gorm:"column:group_name" json:"group_name"`
	AppName        string       `gorm:"column:app_name" json:"app_name"`
	AppAddr        string       `gorm:"column:app_addr" json:"app_addr"`
	CreateTime     str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime     str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*SystemLogRequest) TableName() string {
	return "system_log_requests"
}

// BeforeCreate 创建前
func (t *SystemLogRequest) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *SystemLogRequest) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
