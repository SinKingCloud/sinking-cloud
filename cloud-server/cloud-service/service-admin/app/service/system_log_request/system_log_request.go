package system_log_request

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

// LogRequest 请求日志表
type LogRequest struct {
	Id            string       `json:"id"`
	RequestId     string       `json:"request_id"`
	RequestIp     string       `json:"request_ip"`
	RequestUrl    string       `json:"request_url"`
	RequestMethod string       `json:"request_method"`
	RequestTime   str.DateTime `json:"request_time"`
	ResponseCode  int          `json:"response_code"`
	ResponseTime  int          `json:"response_time"`
	GroupName     string       `json:"group_name"`
	AppName       string       `json:"app_name"`
	AppAddr       string       `json:"app_addr"`
	CreateTime    str.DateTime `json:"create_time"`
	UpdateTime    str.DateTime `json:"update_time"`
}
