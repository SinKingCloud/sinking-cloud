package captcha

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

// Check 判断验证码是否正确
func (*service) Check(key string, code string) bool {
	client, err := rpc.New(group.CloudService, name.ServiceAuth).Call("/rpc/captcha/check", &rpc.Param{
		"key":  key,
		"code": code,
	}).Response()
	if err != nil {
		return false
	}
	return client.Data.(bool)
}
