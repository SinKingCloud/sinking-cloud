package face

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/name"
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

// Check 判断验证码是否正确
func (*Service) Check(id string) (res bool, userId int, cardNo string, cardName string) {
	type Response struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		RequestId string `json:"request_id"`
		Data      *struct {
			UserId   int    `json:"user_id"`
			Result   bool   `json:"result"`
			CardNo   string `json:"card_no"`
			CardName string `json:"card_name"`
		} `json:"data"`
	}
	var response Response
	_, err := rpc.New(group.CloudService, name.ServiceUser).Call("/rpc/face/check", &rpc.Param{
		"id": id,
	}).Unmarshal(&response)
	if err != nil || response.Code != 200 || response.Data == nil {
		return false, 0, "", ""
	}
	return response.Data.Result, response.Data.UserId, response.Data.CardNo, response.Data.CardName
}
