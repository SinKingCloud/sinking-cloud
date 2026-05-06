package aliyun

import (
	"errors"
	"strings"

	"github.com/aliyun/credentials-go/credentials"
)

// Aliyun 聚合阿里云相关工具客户端。
type Aliyun struct {
	Mail    *Mail
	Sms     *Sms
	Company *Company
	Face    *Face
}

// New 创建阿里云工具聚合实例。
func New(accessKeyId, accessKeySecret string) (*Aliyun, error) {
	accessKeyId = strings.TrimSpace(accessKeyId)
	accessKeySecret = strings.TrimSpace(accessKeySecret)
	if accessKeyId == "" {
		return nil, errors.New("AccessKeyId不能为空")
	}
	if accessKeySecret == "" {
		return nil, errors.New("AccessKeySecret不能为空")
	}
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(accessKeyId).
		SetAccessKeySecret(accessKeySecret))
	if err != nil {
		return nil, err
	}
	return &Aliyun{
		Mail:    NewMail(credential),
		Sms:     NewSms(credential),
		Company: NewCompany(credential),
		Face:    NewFace(credential),
	}, nil
}
