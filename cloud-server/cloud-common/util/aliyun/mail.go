package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

var Mail mailIns

type mailIns struct {
}

// Send 发送邮件
// accessKeyId:     阿里云访问密钥ID
// accessKeySecret: 阿里云访问密钥密钥
// accountName:     发件人地址（管理控制台配置的地址）
// nickName:        发信人昵称
// toAddress:       收件人邮箱地址，多个用逗号分隔
// subject:         邮件主题
// htmlBody:        邮件正文（HTML格式）
func (mailIns) Send(accessKeyId, accessKeySecret, accountName, nickName, toAddress, subject, htmlBody string) bool {
	akCredential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(accessKeyId).
		SetAccessKeySecret(accessKeySecret))
	if err != nil {
		return false
	}
	client, err := dm20151123.NewClient(&openapi.Config{
		Credential: akCredential,
		Endpoint:   tea.String("dm.aliyuncs.com"),
	})
	if err != nil {
		return false
	}
	singleSendMailRequest := &dm20151123.SingleSendMailRequest{
		AccountName:    tea.String(accountName),
		AddressType:    tea.Int32(1),
		ReplyToAddress: tea.Bool(true),
		FromAlias:      tea.String(nickName),
		ToAddress:      tea.String(toAddress),
		Subject:        tea.String(subject),
		HtmlBody:       tea.String(htmlBody),
	}
	_, err = client.SingleSendMail(singleSendMailRequest)
	if err != nil {
		return false
	}
	return true
}
