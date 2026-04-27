package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// OrderLog 订单日志表
type OrderLog struct {
	Id           int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	OrderId      int64        `gorm:"column:order_id" json:"order_id"`
	OperatorId   int64        `gorm:"column:operator_id" json:"operator_id"`
	OperatorType int64        `gorm:"column:operator_type" json:"operator_type"`
	Type         int          `gorm:"column:type" json:"type"`
	Content      string       `gorm:"column:content" json:"content"`
	BeforeInfo   string       `gorm:"column:before_info" json:"before_info"`
	AfterInfo    string       `gorm:"column:after_info" json:"after_info"`
	CreateTime   str.DateTime `gorm:"column:create_time" json:"create_time"`
}

// TableName 获取表名
func (*OrderLog) TableName() string {
	return "cloud_order_logs"
}

// BeforeCreate 创建前
func (t *OrderLog) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}
