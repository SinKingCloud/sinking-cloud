package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// File 存储通道
type File struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	ChannelId  int64        `gorm:"column:channel_id" json:"channel_id"`
	SceneId    int64        `gorm:"column:scene_id" json:"scene_id"`
	Name       string       `gorm:"column:name" json:"name"`
	Ext        string       `gorm:"column:ext" json:"ext"`
	Hash       string       `gorm:"column:hash" json:"hash"`
	Size       int64        `gorm:"column:size" json:"size"`
	Path       string       `gorm:"column:path" json:"path"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*File) TableName() string {
	return "cloud_files"
}

// BeforeCreate 创建前
func (t *File) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *File) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
