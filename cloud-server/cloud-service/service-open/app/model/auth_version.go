package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AuthVersion 版本记录
type AuthVersion struct {
	Id             int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	AuthId         int          `gorm:"column:auth_id" json:"auth_id"`
	Version        string       `gorm:"column:version" json:"version"`
	VersionNumber  int          `gorm:"column:version_number" json:"version_number"`
	Title          string       `gorm:"column:title" json:"title"`
	Content        string       `gorm:"column:content" json:"content"`
	ReleaseUrl     string       `gorm:"column:release_url" json:"release_url"`
	ReleaseHash    string       `gorm:"column:release_hash" json:"release_hash"`
	ReleaseDownNum int          `gorm:"column:release_down_num" json:"release_down_num"`
	UpdateUrl      string       `gorm:"column:update_url" json:"update_url"`
	UpdateHash     string       `gorm:"column:update_hash" json:"update_hash"`
	UpdateDownNum  int          `gorm:"column:update_down_num" json:"update_down_num"`
	BuildScript    string       `gorm:"column:build_script" json:"build_script"`
	Status         int          `gorm:"column:status" json:"status"`
	CreateTime     str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime     str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (*AuthVersion) TableName() string {
	return "auth_versions"
}

// BeforeCreate 创建前
func (t *AuthVersion) BeforeCreate(_ *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AuthVersion) BeforeUpdate(_ *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
