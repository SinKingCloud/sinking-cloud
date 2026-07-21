package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// MindVersion 导图版本记录表
type MindVersion struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	DataId     int64        `gorm:"column:data_id" json:"data_id"`
	Type       int          `gorm:"column:type" json:"type"`
	Value      string       `gorm:"column:value" json:"value"`
	Remark     string       `gorm:"column:remark" json:"remark"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*MindVersion) TableName() string {
	return "cloud_mind_versions"
}

// BeforeCreate 创建前
func (t *MindVersion) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *MindVersion) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
