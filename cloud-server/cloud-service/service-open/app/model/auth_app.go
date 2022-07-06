package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AuthApp 程序
type AuthApp struct {
	Id             int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Name           string       `gorm:"column:name" json:"name"`
	Describe       string       `gorm:"column:describe" json:"describe"`
	VersionId      int          `gorm:"column:version_id" json:"version_id"`
	ReleaseDownNum int          `gorm:"column:release_down_num" json:"release_down_num"`
	UpdateDownNum  int          `gorm:"column:update_down_num" json:"update_down_num"`
	Price          float64      `gorm:"column:price" json:"price"`
	Status         int          `gorm:"column:status" json:"status"`
	CreateTime     str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime     str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (AuthApp) TableName() string {
	return "auth_apps"
}

// BeforeCreate 创建前
func (t *AuthApp) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AuthApp) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
