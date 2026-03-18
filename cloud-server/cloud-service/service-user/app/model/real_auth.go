package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// RealAuth 实名记录表
type RealAuth struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int64        `gorm:"column:user_id" json:"user_id"`
	CardName   string       `gorm:"column:card_name" json:"card_name"`
	CardNo     string       `gorm:"column:card_no" json:"card_no"`
	Files      string       `gorm:"column:files" json:"files"`
	Status     int          `gorm:"column:status" json:"status"`
	AuthTime   str.DateTime `gorm:"column:auth_time" json:"auth_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*RealAuth) TableName() string {
	return "cloud_real_auths"
}

// BeforeCreate 创建前
func (t *RealAuth) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *RealAuth) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
