package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// SystemFile 系统文件表
type SystemFile struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Type       int          `gorm:"column:type" json:"type"`
	Hash       string       `gorm:"column:hash" json:"hash"`
	Ext        string       `gorm:"column:ext" json:"ext"`
	Name       string       `gorm:"column:name" json:"name"`
	Path       string       `gorm:"column:path" json:"path"`
	Url        string       `gorm:"column:url" json:"url"`
	Size       int64        `gorm:"column:size" json:"size"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*SystemFile) TableName() string {
	return "system_files"
}

// BeforeCreate 创建前
func (t *SystemFile) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *SystemFile) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
