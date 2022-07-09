package auth_version

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

const (
	Normal = iota //正常
	Error         //封禁
)

// AuthVersion 版本记录
type AuthVersion struct {
	Id             int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	AuthId         int          `gorm:"column:auth_id" json:"auth_id"`
	Version        string       `gorm:"column:version" json:"version"`
	VersionNumber  int          `gorm:"column:version_number" json:"version_number"`
	Title          string       `gorm:"column:title" json:"title"`
	Content        string       `gorm:"column:content" json:"content"`
	ReleaseDownNum int          `gorm:"column:release_down_num" json:"release_down_num"`
	UpdateDownNum  int          `gorm:"column:update_down_num" json:"update_down_num"`
	Status         int          `gorm:"column:status" json:"status"`
	CreateTime     str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime     str.DateTime `gorm:"column:update_time" json:"update_time"`
}
