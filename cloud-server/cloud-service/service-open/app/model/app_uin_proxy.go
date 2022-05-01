package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppUinProxy 代理池列表
type AppUinProxy struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	AreaId     int          `gorm:"column:area_id" json:"area_id"`
	Type       string       `gorm:"column:type" json:"type"`
	ExtractNum int          `gorm:"column:extract_num" json:"extract_num"`
	Ip         string       `gorm:"column:ip" json:"ip"`
	Port       int          `gorm:"column:port" json:"port"`
	Auth       string       `gorm:"column:auth" json:"auth"`
	Status     int          `gorm:"column:status" json:"status"`
	CheckTime  str.DateTime `gorm:"column:check_time" json:"check_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (AppUinProxy) TableName() string {
	return "app_uin_proxys"
}

// BeforeCreate 创建前
func (t *AppUinProxy) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppUinProxy) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
