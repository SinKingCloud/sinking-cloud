package spider

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

// OcrCaptcha 获取验证码cookie
func (*service) OcrCaptcha(url string) string {
	type Response struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		RequestId string `json:"request_id"`
		Data      string `json:"data"`
	}
	var response Response
	_, err := rpc.New(group.CloudService, name.ServiceSpider).Call("/rpc/spider/ocr_captcha", &rpc.Param{
		"captcha_url": url,
	}).Unmarshal(&response)
	if err != nil || response.Code != 200 || response.Data == "" {
		return ""
	}
	return response.Data
}
