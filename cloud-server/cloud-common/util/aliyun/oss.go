package aliyun

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var Oss ossIns

type ossIns struct {
}

func (ossIns) Upload(accessKeyId string, accessKeySecret string, endPoint string, bucketName string, file []byte, path string) bool {
	fmt.Println(endPoint, accessKeyId, accessKeySecret, bucketName, path)
	client, err := oss.New(endPoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("step1", err)
		return false
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("step2", err)
		return false
	}
	err = bucket.PutObject(path, bytes.NewReader(file))
	if err != nil {
		fmt.Println("step3", err)
		return false
	}
	return true
}
