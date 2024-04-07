package system_log_request

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

// LogRequest 请求日志表
type LogRequest struct {
	Id            string       `gorm:"column:id;PRIMARY_KEY" json:"id"`
	RequestId     string       `gorm:"column:request_id" json:"request_id"`
	RequestIp     string       `gorm:"column:request_ip" json:"request_ip"`
	RequestUrl    string       `gorm:"column:request_url" json:"request_url"`
	RequestMethod string       `gorm:"column:request_method" json:"request_method"`
	RequestTime   str.DateTime `gorm:"column:request_time" json:"request_time"`
	ResponseCode  int          `gorm:"column:response_code" json:"response_code"`
	ResponseTime  int          `gorm:"column:response_time" json:"response_time"`
	GroupName     string       `gorm:"column:group_name" json:"group_name"`
	AppName       string       `gorm:"column:app_name" json:"app_name"`
	AppAddr       string       `gorm:"column:app_addr" json:"app_addr"`
	CreateTime    str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime    str.DateTime `gorm:"column:update_time" json:"update_time"`
}
