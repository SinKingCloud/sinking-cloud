package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

var Sms sms

type sms struct {
}

// Send 发送短信
// accessKeyId:     阿里云访问密钥ID
// accessKeySecret: 阿里云访问密钥密钥
// phone:           接收手机号
// signName:        短信签名
// templateCode:    短信模板CODE
// templateParam:   短信模板参数，JSON格式，如：{\"code\":\"1234\"}
func (sms) Send(accessKeyId, accessKeySecret, phone, signName, templateCode, templateParam string) bool {
	akCredential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(accessKeyId).
		SetAccessKeySecret(accessKeySecret))
	if err != nil {
		return false
	}
	client, err := dysmsapi20170525.NewClient(&openapi.Config{
		Credential: akCredential,
		Endpoint:   tea.String("dysmsapi.aliyuncs.com"),
	})
	if err != nil {
		return false
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		PhoneNumbers:  tea.String(phone),
		TemplateParam: tea.String(templateParam),
	}
	result, err := client.SendSms(sendSmsRequest)
	if err != nil {
		return false
	}
	if result.Body.Code != nil && *result.Body.Code == "OK" {
		return true
	}
	return false
}
