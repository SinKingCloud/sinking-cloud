package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// Scene 业务场景
type Scene struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Name       string       `gorm:"column:name" json:"name"`
	MaxSize    int64        `gorm:"column:max_size" json:"max_size"`
	AllowExt   string       `gorm:"column:allow_ext" json:"allow_ext"`
	Login      int          `gorm:"column:login" json:"login"`
	ChannelId  int64        `gorm:"column:channel_id" json:"channel_id"`
	Remark     string       `gorm:"column:remark" json:"remark"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*Scene) TableName() string {
	return "cloud_scenes"
}

// BeforeCreate 创建前
func (t *Scene) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *Scene) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
