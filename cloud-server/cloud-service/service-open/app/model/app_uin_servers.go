package model

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"gorm.io/gorm"
	"time"
)

// AppUinServer QQ服务器列表
type AppUinServer struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Key        string       `gorm:"column:key" json:"key"`
	Name       string       `gorm:"column:name" json:"name"`
	Addr       string       `gorm:"column:addr" json:"addr"`
	Proxy      string       `gorm:"column:proxy" json:"proxy"`
	UinNum     int          `gorm:"column:uin_num" json:"uin_num"`
	Status     int          `gorm:"column:status" json:"status"`
	CheckTime  str.DateTime `gorm:"column:check_time" json:"check_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// TableName 获取表名
func (AppUinServer) TableName() string {
	return "app_uin_servers"
}

// BeforeCreate 创建前
func (t *AppUinServer) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = str.DateTime(time.Now())
	return nil
}

// BeforeUpdate 更新前
func (t *AppUinServer) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = str.DateTime(time.Now())
	return nil
}
