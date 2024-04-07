package rpc

import (
	"encoding/json"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	util2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-go/sinking-sdk/sinking-sdk-go"
)

type Rpc struct {
	groupName string
	name      string
	res       string
	err       error
	requestId string
}

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

type Param map[string]interface{}

// New 新建调用实例
func New(groupName string, name string) *Rpc {
	return &Rpc{
		groupName: groupName,
		name:      name,
	}
}

// Call 调用指定地址
func (r *Rpc) Call(url string, param *Param) *Rpc {
	r.requestId = str.GetSnowWorkIns().GetUuid()
	r.res, r.err = util.Consul.Rpc(r.groupName, r.name).Method("POST").Header(map[string]string{
		util2.RequestId: r.requestId,
		util.Conf.GetString(config.ConfRegisterTokenKey): util.Conf.GetString(config.ConfRegisterTokenValue),
		util2.RpcHeaderTokenName:                         util.Conf.GetString(config.ConfRegisterTokenValue),
	}).Timeout(60).ReTry(5).Call(url, &sinking_sdk_go.Param{
		"data":  param,
		"group": util.Conf.GetString(config.ConfApplicationGroup),
		"name":  util.Conf.GetString(config.ConfApplicationName),
	}, sinking_sdk_go.Poll)
	return r
}

// Body 获取结果
func (r *Rpc) Body() (string, error) {
	return r.res, r.err
}

// Unmarshal 获取结果
func (r *Rpc) Unmarshal(v interface{}) (interface{}, error) {
	if r.err != nil {
		return nil, r.err
	}
	err := json.Unmarshal([]byte(r.res), v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Response 获取响应对象
func (r *Rpc) Response() (*Response, error) {
	var response Response
	data, err := r.Unmarshal(&response)
	if err != nil {
		return nil, err
	}
	return data.(*Response), nil
}
