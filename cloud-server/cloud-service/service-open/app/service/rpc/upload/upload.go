package upload

import (
	"encoding/base64"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/name"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/rpc"
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

// Temp 上传至缓存文件夹
func (*service) Temp(fileName string, data []byte) string {
	type Response struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		RequestId string `json:"request_id"`
		Data      string `json:"data"`
	}
	var response Response
	_, err := rpc.New(group.CloudService, name.ServiceAuth).Call("/rpc/upload/temp", &rpc.Param{
		"name": fileName,
		"data": base64.StdEncoding.EncodeToString(data),
	}).Unmarshal(&response)
	if err != nil || response.Code != 200 || response.Data == "" {
		return ""
	}
	return response.Data
}
