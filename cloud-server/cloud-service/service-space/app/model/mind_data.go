package model

import (
	"time"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
)

// MindData 导图数据表
type MindData struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int64        `gorm:"column:user_id" json:"user_id"`
	TypeId     int64        `gorm:"column:type_id" json:"type_id"`
	Title      string       `gorm:"column:title" json:"title"`
	Describe   string       `gorm:"column:describe" json:"describe"`
	Value      string       `gorm:"column:value" json:"value"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*MindData) TableName() string {
	return "cloud_mind_datas"
}

// BeforeCreate 创建前
func (t *MindData) BeforeCreate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *MindData) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
