package notice

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"sync"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// WebNotice 网站通知
type WebNotice struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Place      string       `gorm:"column:place" json:"place"`
	Title      string       `gorm:"column:title" json:"title"`
	Status     int          `gorm:"column:status" json:"status"`
	Sort       int          `gorm:"column:sort" json:"sort"`
	LookNum    int          `gorm:"column:look_num" json:"look_num"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}
