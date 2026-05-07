package aliyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	cloudauth20190307 "github.com/alibabacloud-go/cloudauth-20190307/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

type Face struct {
	Credential credentials.Credential
}

// FaceInitResult 是初始化人脸实名认证后的认证信息。
type FaceInitResult struct {
	CertifyId  string `json:"certify_id"`  // 阿里云认证ID，后续查询认证结果时使用。
	CertifyUrl string `json:"certify_url"` // 用户认证页面地址。
}

// FaceQueryResult 是人脸实名认证查询结果。
type FaceQueryResult struct {
	Passed    bool     `json:"passed"`      // 是否认证通过
	Reason    string   `json:"reason"`      // 认证结果说明
	Bucket    string   `json:"bucket_name"` // 资料储存bucket名称
	PicUrls   []string `json:"pic_urls"`    // 认证照片外链
	PicPaths  []string `json:"pic_paths"`   // 认证照片路径
	VideoUrl  string   `json:"video_url"`   // 认证视频外链
	VideoPath string   `json:"video_path"`  // 认证视频路径
}

// NewFace 创建阿里云人脸实名认证客户端实例。
func NewFace(credential credentials.Credential) *Face {
	return &Face{Credential: credential}
}

// Init 初始化人脸实名认证，返回认证ID和认证地址。
// sceneId:      阿里云认证场景ID
// outerOrderNo: 商户侧自定义业务唯一标识
// cardNo:       身份证号
// cardName:     姓名
// returnUrl:    认证完成后跳转地址
// metaInfo:     设备环境信息
func (ins *Face) Init(sceneId int64, outerOrderNo, cardNo, cardName, returnUrl, callbackUrl, metaInfo string) (*FaceInitResult, error) {
	cardNo = strings.TrimSpace(cardNo)
	outerOrderNo = strings.TrimSpace(outerOrderNo)
	cardName = strings.TrimSpace(cardName)
	returnUrl = strings.TrimSpace(returnUrl)
	callbackUrl = strings.TrimSpace(callbackUrl)
	metaInfo = strings.TrimSpace(metaInfo)
	if sceneId <= 0 {
		return nil, errors.New("认证场景ID不能为空")
	}
	if outerOrderNo == "" {
		return nil, errors.New("商户业务唯一标识不能为空")
	}
	if cardNo == "" {
		return nil, errors.New("身份证号不能为空")
	}
	if cardName == "" {
		return nil, errors.New("姓名不能为空")
	}
	if metaInfo == "" {
		return nil, errors.New("设备环境信息不能为空")
	}
	if ins == nil || ins.Credential == nil {
		return nil, errors.New("阿里云凭据不能为空")
	}
	client, err := cloudauth20190307.NewClient(&openapi.Config{
		Credential: ins.Credential,
		Endpoint:   tea.String("cloudauth.aliyuncs.com"),
	})
	if err != nil {
		return nil, fmt.Errorf("创建阿里云实人认证客户端失败:%w", err)
	}
	request := &cloudauth20190307.InitFaceVerifyRequest{
		CertNo:            tea.String(cardNo),
		OuterOrderNo:      tea.String(outerOrderNo),
		SceneId:           tea.Int64(sceneId),
		ProductCode:       tea.String("ID_PRO"),
		CertType:          tea.String("IDENTITY_CARD"),
		CertName:          tea.String(cardName),
		ReturnUrl:         tea.String(returnUrl),
		CallbackUrl:       tea.String(callbackUrl),
		MetaInfo:          tea.String(metaInfo),
		ProcedurePriority: tea.String("url"),
		VideoEvidence:     tea.String("true"),
	}
	resp, err := client.InitFaceVerify(request)
	if err != nil {
		return nil, ins.formatSdkError("初始化人脸实名认证失败", err)
	}
	if resp == nil || resp.Body == nil {
		return nil, errors.New("初始化人脸实名认证失败:响应为空")
	}
	body := resp.Body
	if tea.StringValue(body.Code) != "200" {
		return nil, fmt.Errorf("初始化人脸实名认证失败:%s", ins.joinMessage(tea.StringValue(body.Code), tea.StringValue(body.Message), tea.StringValue(body.RequestId)))
	}
	if body.ResultObject == nil {
		return nil, errors.New("初始化人脸实名认证失败:认证结果为空")
	}
	result := &FaceInitResult{
		CertifyId:  tea.StringValue(body.ResultObject.CertifyId),
		CertifyUrl: tea.StringValue(body.ResultObject.CertifyUrl),
	}
	if result.CertifyId == "" {
		return nil, errors.New("初始化人脸实名认证失败:认证ID为空")
	}
	if result.CertifyUrl == "" {
		return nil, errors.New("初始化人脸实名认证失败:认证地址为空")
	}
	return result, nil
}

// Query 查询人脸实名认证结果。
// certifyId: Init返回的认证ID
// sceneId:   阿里云认证场景ID
func (ins *Face) Query(certifyId string, sceneId int64) (*FaceQueryResult, error) {
	certifyId = strings.TrimSpace(certifyId)
	if certifyId == "" {
		return nil, errors.New("认证ID不能为空")
	}
	if sceneId <= 0 {
		return nil, errors.New("认证场景ID不能为空")
	}
	if ins == nil || ins.Credential == nil {
		return nil, errors.New("阿里云凭据不能为空")
	}
	client, err := cloudauth20190307.NewClient(&openapi.Config{
		Credential: ins.Credential,
		Endpoint:   tea.String("cloudauth.aliyuncs.com"),
	})
	if err != nil {
		return nil, fmt.Errorf("创建阿里云实人认证客户端失败:%w", err)
	}
	request := &cloudauth20190307.DescribeFaceVerifyRequest{
		CertifyId: tea.String(certifyId),
		SceneId:   tea.Int64(sceneId),
	}
	resp, err := client.DescribeFaceVerify(request)
	if err != nil {
		return nil, ins.formatSdkError("查询人脸实名认证结果失败", err)
	}
	if resp == nil || resp.Body == nil {
		return nil, errors.New("查询人脸实名认证结果失败:响应为空")
	}
	body := resp.Body
	if tea.StringValue(body.Code) != "200" {
		return nil, fmt.Errorf("查询人脸实名认证结果失败:%s", ins.joinMessage(tea.StringValue(body.Code), tea.StringValue(body.Message), tea.StringValue(body.RequestId)))
	}
	if body.ResultObject == nil {
		return nil, errors.New("查询人脸实名认证结果失败:认证结果为空")
	}
	resultObject := body.ResultObject
	subCode := tea.StringValue(resultObject.SubCode)
	passedFlag := tea.StringValue(resultObject.Passed)
	result := &FaceQueryResult{
		Passed: passedFlag == "T" || subCode == "200",
		Reason: ins.subCodeText(subCode),
	}
	ins.fillMaterialInfo(result, tea.StringValue(resultObject.MaterialInfo))
	if !result.Passed {
		if result.Reason != "" {
			return result, fmt.Errorf("人脸实名认证未通过:%s", result.Reason)
		}
		return result, errors.New("人脸实名认证未通过")
	}
	return result, nil
}

// fillMaterialInfo 将阿里云MaterialInfo解析成业务易用字段。
func (ins *Face) fillMaterialInfo(result *FaceQueryResult, materialInfo string) {
	if result == nil || strings.TrimSpace(materialInfo) == "" {
		return
	}
	type facePictureInfo struct {
		CustomFaceObjectNameUrls []string `json:"customFaceObjectNameUrls"`
		CustomFaceObjectNames    []string `json:"customFaceObjectNames"`
		OssBucketName            string   `json:"ossBucketName"`
		OssObjectName            string   `json:"ossObjectName"`
		OssVerifyVideoObjectName string   `json:"ossVerifyVideoObjectName"`
		PictureUrl               string   `json:"pictureUrl"`
		VerifyVideoUrl           string   `json:"verifyVideoUrl"`
	}
	type faceMaterialInfo struct {
		FacialPictureFront facePictureInfo `json:"facialPictureFront"`
	}
	var info faceMaterialInfo
	if err := json.NewDecoder(strings.NewReader(materialInfo)).Decode(&info); err != nil {
		return
	}
	pictureInfo := info.FacialPictureFront
	result.Bucket = pictureInfo.OssBucketName
	result.PicPaths = append(result.PicPaths, pictureInfo.CustomFaceObjectNames...)
	result.PicUrls = append(result.PicUrls, pictureInfo.CustomFaceObjectNameUrls...)
	if pictureInfo.OssObjectName != "" {
		result.PicPaths = append(result.PicPaths, pictureInfo.OssObjectName)
	}
	if pictureInfo.PictureUrl != "" {
		result.PicUrls = append(result.PicUrls, pictureInfo.PictureUrl)
	}
	result.VideoPath = pictureInfo.OssVerifyVideoObjectName
	result.VideoUrl = pictureInfo.VerifyVideoUrl
}

// formatSdkError 将阿里云SDK错误转换成业务可展示的错误信息。
func (ins *Face) formatSdkError(prefix string, err error) error {
	if err == nil {
		return nil
	}
	var sdkErr *tea.SDKError
	if !errors.As(err, &sdkErr) {
		return fmt.Errorf("%s:%w", prefix, err)
	}
	message := tea.StringValue(sdkErr.Message)
	if message == "" {
		message = err.Error()
	}
	recommend := ins.getErrorRecommend(tea.StringValue(sdkErr.Data))
	if recommend != "" {
		return fmt.Errorf("%s:%s,建议:%s", prefix, message, recommend)
	}
	return fmt.Errorf("%s:%s", prefix, message)
}

// getErrorRecommend 从阿里云SDK错误扩展数据中提取诊断建议。
func (ins *Face) getErrorRecommend(data string) string {
	if data == "" {
		return ""
	}
	var result map[string]interface{}
	if err := json.NewDecoder(strings.NewReader(data)).Decode(&result); err != nil {
		return ""
	}
	recommend, ok := result["Recommend"]
	if !ok || recommend == nil {
		return ""
	}
	return fmt.Sprint(recommend)
}

// joinMessage 拼接阿里云非成功响应的错误信息。
func (ins *Face) joinMessage(code, message, requestId string) string {
	items := make([]string, 0, 3)
	if code != "" {
		items = append(items, "Code:"+code)
	}
	if message != "" {
		items = append(items, "Message:"+message)
	}
	if requestId != "" {
		items = append(items, "RequestId:"+requestId)
	}
	if len(items) == 0 {
		return "未知错误"
	}
	return strings.Join(items, ",")
}

// subCodeText 将阿里云认证结果码转换成中文说明。
func (ins *Face) subCodeText(subCode string) string {
	switch strings.TrimSpace(subCode) {
	case "200":
		return "认证通过"
	case "201":
		return "姓名和身份证不一致"
	case "202":
		return "查询不到身份信息"
	case "203":
		return "查询不到照片或照片不可用"
	case "204":
		return "人脸比对不一致"
	case "205":
		return "活体检测存在风险"
	case "206":
		return "业务策略限制"
	default:
		if subCode == "" {
			return ""
		}
		return "未知认证结果码:" + subCode
	}
}
