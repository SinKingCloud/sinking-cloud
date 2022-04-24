package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppProxyArea 代理池地区列表
type AppProxyArea struct {
	Id          int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Province    string       `gorm:"column:province" json:"province"`
	ProvinceKey string       `gorm:"column:province_key" json:"province_key"`
	City        string       `gorm:"column:city" json:"city"`
	CityKey     string       `gorm:"column:city_key" json:"city_key"`
	Status      int          `gorm:"column:status" json:"status"`
	CheckTime   str.DateTime `gorm:"column:check_time" json:"check_time"`
	CreateTime  str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime  str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (AppProxyArea) TableName() string {
	return "app_proxy_areas"
}

// BeforeCreate 创建前
func (t *AppProxyArea) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppProxyArea) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
