package user_app

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/name"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"sync"
)

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

type UserApp struct {
	Id         int          `json:"id"`
	UserId     int          `json:"user_id"`
	Key        string       `json:"key"`
	Name       string       `json:"name"`
	Describe   string       `json:"describe"`
	Status     int          `json:"status"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}

// FindById 获取用户信息
func (*Service) FindById(id int) *UserApp {
	type Response struct {
		Code      int      `json:"code"`
		Message   string   `json:"message"`
		RequestId string   `json:"request_id"`
		Data      *UserApp `json:"data"`
	}
	var response Response
	_, err := rpc.New(group.CloudService, name.ServiceOpen).Call("/rpc/user_app/findById", &rpc.Param{
		"id": id,
	}).Unmarshal(&response)
	if err != nil || response.Code != 200 || response.Data == nil {
		return nil
	}
	return response.Data
}
