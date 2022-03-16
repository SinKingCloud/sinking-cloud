package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

var Sms smsIns

type smsIns struct {
}

func (smsIns) Send(accessKeyId string, accessKeySecret string, phone string, signName string, templateCode string, templateParam string) bool {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = signName
	request.TemplateCode = templateCode
	request.TemplateParam = templateParam
	response, err := client.SendSms(request)
	if err == nil {
		if response.Message == "OK" {
			return true
		}
		return false
	}
	return false
}
