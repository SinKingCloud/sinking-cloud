package notice

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"sync"
)

// service 单例对象
type service struct {
}

// obj 单例对象
var (
	obj  *service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *service {
	once.Do(func() {
		obj = &service{}
	})
	return obj
}

// WebNotice 网站通知
type WebNotice struct {
	Id         int          `json:"id"`
	Place      string       `json:"place"`
	Title      string       `json:"title"`
	Status     int          `json:"status"`
	Sort       int          `json:"sort"`
	LookNum    int          `json:"look_num"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}
