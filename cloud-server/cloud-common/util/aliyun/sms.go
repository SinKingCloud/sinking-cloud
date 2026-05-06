package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

type Sms struct {
	Credential credentials.Credential
}

// NewSms 创建阿里云短信客户端实例。
func NewSms(credential credentials.Credential) *Sms {
	return &Sms{Credential: credential}
}

// Send 发送短信。
// phone:         接收手机号
// signName:      短信签名
// templateCode:  短信模板CODE
// templateParam: 短信模板参数，JSON格式，如：{\"code\":\"1234\"}
func (ins *Sms) Send(phone, signName, templateCode, templateParam string) bool {
	if ins == nil || ins.Credential == nil {
		return false
	}
	client, err := dysmsapi20170525.NewClient(&openapi.Config{
		Credential: ins.Credential,
		Endpoint:   tea.String("dysmsapi.aliyuncs.com"),
	})
	if err != nil {
		return false
	}
	request := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(templateCode),
		PhoneNumbers:  tea.String(phone),
		TemplateParam: tea.String(templateParam),
	}
	result, err := client.SendSms(request)
	if err != nil {
		return false
	}
	return result.Body.Code != nil && *result.Body.Code == "OK"
}
