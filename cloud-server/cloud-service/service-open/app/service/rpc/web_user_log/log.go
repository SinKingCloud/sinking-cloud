package web_user_log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/name"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
)

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

// Create 创建用户日志
func (*service) Create(userId int, requestId string, requestIp string, types int, title string, content string) {
	_, _ = rpc.New(group.CloudService, name.ServiceUser).Call("/rpc/log/create", &rpc.Param{
		"user_id":    userId,
		"request_id": requestId,
		"request_ip": requestIp,
		"type":       types,
		"title":      title,
		"content":    content,
	}).Response()
}
