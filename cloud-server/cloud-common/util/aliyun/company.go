package aliyun

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	cloudauth20221125 "github.com/alibabacloud-go/cloudauth-20221125/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

type Company struct {
	Credential credentials.Credential
}

// NewCompany 创建阿里云企业核验客户端实例。
func NewCompany(credential credentials.Credential) *Company {
	return &Company{Credential: credential}
}

// Verify 企业四要素核验。
// companyName: 企业名称
// companyNo:   统一社会信用代码
// cardName:    法人姓名
// cardNo:      法人身份证号
func (ins *Company) Verify(companyName, companyNo, cardName, cardNo string) error {
	companyName = strings.TrimSpace(companyName)
	companyNo = strings.TrimSpace(companyNo)
	cardName = strings.TrimSpace(cardName)
	cardNo = strings.TrimSpace(cardNo)
	if companyName == "" {
		return errors.New("企业名称不能为空")
	}
	if companyNo == "" {
		return errors.New("统一社会信用代码不能为空")
	}
	if cardName == "" {
		return errors.New("法人姓名不能为空")
	}
	if cardNo == "" {
		return errors.New("法人身份证号不能为空")
	}
	if ins == nil || ins.Credential == nil {
		return errors.New("阿里云凭据不能为空")
	}
	client, err := cloudauth20221125.NewClient(&openapi.Config{
		Credential: ins.Credential,
		Endpoint:   tea.String("cloudauth.aliyuncs.com"),
	})
	if err != nil {
		return fmt.Errorf("创建阿里云企业核验客户端失败:%w", err)
	}
	request := &cloudauth20221125.EntElementVerifyV2Request{
		SceneCode:         tea.String("1000000006"),
		MerchantBizId:     tea.String(ins.newMerchantBizId()),
		MerchantUserId:    tea.String(companyNo),
		UserAuthorization: tea.String("1"),
		InfoVerifyType:    tea.String("ENT_4META"),
		LicenseNo:         tea.String(companyNo),
		EntName:           tea.String(companyName),
		LegalPersonName:   tea.String(cardName),
		LegalPersonCertNo: tea.String(cardNo),
	}
	resp, err := client.EntElementVerifyV2(request)
	if err != nil {
		return ins.formatSdkError(err)
	}
	if resp == nil || resp.Body == nil {
		return errors.New("企业四要素核验失败:响应为空")
	}
	body := resp.Body
	if tea.StringValue(body.Code) != "Success" {
		return fmt.Errorf("企业四要素核验失败:%s", ins.joinMessage(tea.StringValue(body.Code), tea.StringValue(body.Message), tea.StringValue(body.RequestId)))
	}
	if body.Result == nil {
		return errors.New("企业四要素核验失败:核验结果为空")
	}
	result := body.Result
	switch tea.StringValue(result.BizCode) {
	case "1":
		return nil
	case "2":
		return fmt.Errorf("企业四要素核验不一致:%s", ins.reasonText(tea.StringValue(result.ReasonDetail)))
	case "3":
		return errors.New("企业四要素核验未查得")
	default:
		return fmt.Errorf("企业四要素核验失败:未知核验结果码%s", tea.StringValue(result.BizCode))
	}
}

// newMerchantBizId 生成阿里云商户侧唯一业务标识。
func (ins *Company) newMerchantBizId() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err == nil {
		return hex.EncodeToString(buf)
	}
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// formatSdkError 将阿里云SDK错误转换成业务可展示的错误信息。
func (ins *Company) formatSdkError(err error) error {
	if err == nil {
		return nil
	}
	var sdkErr *tea.SDKError
	if !errors.As(err, &sdkErr) {
		return fmt.Errorf("企业四要素核验失败:%w", err)
	}
	message := tea.StringValue(sdkErr.Message)
	if message == "" {
		message = err.Error()
	}
	recommend := ins.getErrorRecommend(tea.StringValue(sdkErr.Data))
	if recommend != "" {
		return fmt.Errorf("企业四要素核验失败:%s,建议:%s", message, recommend)
	}
	return fmt.Errorf("企业四要素核验失败:%s", message)
}

// getErrorRecommend 从阿里云SDK错误扩展数据中提取诊断建议。
func (ins *Company) getErrorRecommend(data string) string {
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
func (ins *Company) joinMessage(code, message, requestId string) string {
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

// reasonText 将阿里云不一致字段标识转换成中文原因。
func (ins *Company) reasonText(reasonDetail string) string {
	reasonDetail = strings.TrimSpace(reasonDetail)
	if reasonDetail == "" {
		return "信息不一致"
	}
	reasonMap := map[string]string{
		"LegalPersonNameFlag":   "法人名称不一致",
		"LegalPersonCertNoFlag": "法人证件号码不一致",
		"EntNameFlag":           "企业名称不一致",
		"LicenseNoFlag":         "营业执照号不一致",
	}
	reasons := make([]string, 0)
	for _, reason := range strings.Split(reasonDetail, ",") {
		reason = strings.TrimSpace(reason)
		if reason == "" {
			continue
		}
		if text, ok := reasonMap[reason]; ok {
			reasons = append(reasons, text)
			continue
		}
		reasons = append(reasons, reason)
	}
	if len(reasons) == 0 {
		return "信息不一致"
	}
	return strings.Join(reasons, ",")
}
