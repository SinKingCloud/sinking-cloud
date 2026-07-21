package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// MindShare 导图分享表
type MindShare struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	DataId     int64        `gorm:"column:data_id" json:"data_id"`
	Code       string       `gorm:"column:code" json:"code"`
	Status     int          `gorm:"column:status" json:"status"`
	ExpireTime str.DateTime `gorm:"column:expire_time" json:"expire_time"`
	Remark     string       `gorm:"column:remark" json:"remark"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*MindShare) TableName() string {
	return "cloud_mind_shares"
}

// BeforeCreate 创建前
func (t *MindShare) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *MindShare) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
