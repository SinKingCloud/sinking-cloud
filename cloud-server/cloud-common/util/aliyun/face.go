package aliyun

import (
	"encoding/json"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

var Face faceIns

type faceIns struct {
	accessKeyId     string
	accessKeySecret string
}

func NewFace(accessKeyId string, accessKeySecret string) *faceIns {
	return &faceIns{accessKeyId: accessKeyId, accessKeySecret: accessKeySecret}
}

type responseContent struct {
	Data struct {
		CertifyId  string `json:"certifyId"`
		CertifyUrl string `json:"certifyUrl"`
		Passed     string `json:"passed"`
	} `json:"data"`
	RequestID string `json:"requestId"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
}

// Init 初始化认证接口
func (f *faceIns) Init(uuid string, sceneId string, bizCode string, certType string, certNo string, certName string, returnUrl string, callbackUrl string) (id string, url string, err error) {
	// 需要替换AccessKey ID和AccessKey Secret。
	client, err := sdk.NewClientWithAccessKey(
		"cn-shanghai", f.accessKeyId, f.accessKeySecret)
	if err != nil {
		return "", "", err
	}
	request := requests.NewCommonRequest()
	// 下述参数无需调整。
	request.Domain = "saf.cn-shanghai.aliyuncs.com"
	request.Version = "2017-03-31"
	request.ApiName = "ExecuteRequest"
	request.Scheme = "https"
	request.QueryParams["Service"] = "fin_face_verify"
	// 下述参数依照实际使用的接口做调整，本演示仅针对认证初始化接口。
	serviceParam := map[string]string{
		"method":       "init",
		"sceneId":      sceneId,
		"outerOrderNo": uuid,
		"bizCode":      bizCode,
		"identityType": "CERT_INFO",
		"certType":     certType,
		"certNo":       certNo,
		"certName":     certName,
		"returnUrl":    returnUrl,
		"callbackUrl":  callbackUrl,
	}
	serviceParamByte, err := json.Marshal(serviceParam)
	if err != nil {
		return "", "", err
	}
	request.QueryParams["ServiceParameters"] = string(serviceParamByte)
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return "", "", err
	}
	respContent := responseContent{}
	err = json.Unmarshal(response.GetHttpContentBytes(), &respContent)
	if err != nil {
		return "", "", err
	}
	if respContent.Message != "OK" {
		return "", "", errors.New(respContent.Message)
	}
	return respContent.Data.CertifyId, respContent.Data.CertifyUrl, nil
}

func (f *faceIns) Query(sceneId string, certifyId string) (bool, error) {
	// 需要替换AccessKey ID和AccessKey Secret。
	client, err := sdk.NewClientWithAccessKey(
		"cn-shanghai", f.accessKeyId, f.accessKeySecret)
	if err != nil {
		return false, err
	}
	request := requests.NewCommonRequest()
	// 下述参数无需调整。
	request.Domain = "saf.cn-shanghai.aliyuncs.com"
	request.Version = "2017-03-31"
	request.ApiName = "ExecuteRequest"
	request.Scheme = "https"
	request.QueryParams["Service"] = "fin_face_verify"
	serviceParam := map[string]string{
		"method":    "query",
		"sceneId":   sceneId,
		"certifyId": certifyId,
	}
	serviceParamByte, err := json.Marshal(serviceParam)
	if err != nil {
		return false, err
	}
	request.QueryParams["ServiceParameters"] = string(serviceParamByte)
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return false, err
	}
	respContent := responseContent{}
	err = json.Unmarshal(response.GetHttpContentBytes(), &respContent)
	if err != nil {
		return false, err
	}
	if respContent.Message != "OK" {
		return false, errors.New(respContent.Message)
	}
	return respContent.Data.Passed == "T", nil
}
