package web_user

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
func (*service) FindById(userId int) *model.WebUser {
	type Response struct {
		Code      int            `json:"code"`
		Message   string         `json:"message"`
		RequestId string         `json:"request_id"`
		Data      *model.WebUser `json:"data"`
	}
	var response Response
	_, err := rpc.New(group.CloudService, name.ServiceUser).Call("/rpc/user/findById", &rpc.Param{
		"user_id": userId,
	}).Unmarshal(&response)
	if err != nil || response.Code != 200 || response.Data == nil {
		return nil
	}
	return response.Data
}

// Money 用户余额操作
func (*service) Money(userId int, mod int, money string) bool {
	response, err := rpc.New(group.CloudService, name.ServiceUser).Call("/rpc/user/money", &rpc.Param{
		"user_id": userId,
		"mod":     mod,
		"money":   money,
	}).Response()
	if err != nil || response.Code != 200 {
		return false
	}
	return true
}
