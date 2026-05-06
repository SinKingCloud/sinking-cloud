package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

type Mail struct {
	Credential credentials.Credential
}

// NewMail 创建阿里云邮件客户端实例。
func NewMail(credential credentials.Credential) *Mail {
	return &Mail{Credential: credential}
}

// Send 发送邮件。
// accountName: 发件人地址（管理控制台配置的地址）
// nickName:    发信人昵称
// toAddress:   收件人邮箱地址，多个用逗号分隔
// subject:     邮件主题
// htmlBody:    邮件正文（HTML格式）
func (ins *Mail) Send(accountName, nickName, toAddress, subject, htmlBody string) bool {
	if ins == nil || ins.Credential == nil {
		return false
	}
	client, err := dm20151123.NewClient(&openapi.Config{
		Credential: ins.Credential,
		Endpoint:   tea.String("dm.aliyuncs.com"),
	})
	if err != nil {
		return false
	}
	request := &dm20151123.SingleSendMailRequest{
		AccountName:    tea.String(accountName),
		AddressType:    tea.Int32(1),
		ReplyToAddress: tea.Bool(true),
		FromAlias:      tea.String(nickName),
		ToAddress:      tea.String(toAddress),
		Subject:        tea.String(subject),
		HtmlBody:       tea.String(htmlBody),
	}
	_, err = client.SingleSendMail(request)
	return err == nil
}
