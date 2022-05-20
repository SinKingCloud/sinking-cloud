package system_user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/name"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
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

// FindById 获取用户信息
func (*service) FindById(userId int) *model.SystemUser {
	type Response struct {
		Code      int               `json:"code"`
		Message   string            `json:"message"`
		RequestId string            `json:"request_id"`
		Data      *model.SystemUser `json:"data"`
	}
	var response Response
	_, err := rpc.New(group.CloudService, name.ServiceAdmin).Call("/rpc/user/findById", &rpc.Param{
		"user_id": userId,
	}).Unmarshal(&response)
	if err != nil || response.Code != 200 || response.Data == nil {
		return nil
	}
	return response.Data
}
