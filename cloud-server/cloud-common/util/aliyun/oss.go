package aliyun

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var Oss ossIns

type ossIns struct {
}

func (ossIns) Upload(accessKeyId string, accessKeySecret string, endPoint string, bucketName string, file []byte, path string) bool {
	client, err := oss.New(endPoint, accessKeyId, accessKeySecret)
	if err != nil {
		return false
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return false
	}
	err = bucket.PutObject(path, bytes.NewReader(file))
	if err != nil {
		return false
	}
	return true
}
