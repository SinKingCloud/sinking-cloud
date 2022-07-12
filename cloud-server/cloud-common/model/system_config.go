package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// SystemConfig 系统配置表
type SystemConfig struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Pid        int          `gorm:"column:pid" json:"pid"`
	Name       string       `gorm:"column:name" json:"name"`
	Describe   string       `gorm:"column:describe" json:"describe"`
	Key        string       `gorm:"column:key" json:"key"`
	Value      string       `gorm:"column:value" json:"value"`
	Control    string       `gorm:"column:control" json:"control"`
	Sort       int          `gorm:"column:sort" json:"sort"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (SystemConfig) TableName() string {
	return "system_configs"
}

// BeforeCreate 创建前
func (t *SystemConfig) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *SystemConfig) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
