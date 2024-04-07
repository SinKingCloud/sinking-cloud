package system_user

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/name"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
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

// FindById 获取用户信息
func (*Service) FindById(userId int) *model.SystemUser {
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
