package uin_tool

import (
	"encoding/json"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	server2 "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_uin/server"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// add 添加账号
func (s *speed) add(uin string, pwd string, phone string, phoneArea string, proto int) bool {
	//判断单线程是否上锁
	if server2.GetIns().ServerIsLock("add_" + s.addr) {
		//有锁的话等待执行完成
		server2.GetIns().ServerWaitLock("add_" + s.addr)
	}
	//加锁
	server2.GetIns().ServerLock("add_" + s.addr)
	defer server2.GetIns().ServerUnLock("add_" + s.addr)
	//阻止全局登陆
	server2.GetIns().ServerInc("add_"+s.addr, 120)
	defer server2.GetIns().ServerDec("add_"+s.addr, 120)
	//等待计数器完成
	server2.GetIns().ServerNumWaitZero("login_" + s.addr)
	p := "0"
	if proto == 1 {
		p = "1"
	}
	url := "/?&type=Android/&pe=into/"
	data := map[string]string{
		"Uin":      uin,
		"Pass":     pwd,
		"brand":    "",
		"model":    "",
		"areacode": phoneArea,
		"protocol": p,
		"phone":    phone,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	if res.State == "0" {
		return true
	}
	time.Sleep(3 * time.Second)
	return false
}

// SmsVerify 短信验证
type SmsVerify struct {
	Retcode   int    `json:"retcode"`
	Retmsg    string `json:"retmsg"`
	NeedCheck bool   `json:"needCheck"`
	AreaCode  string `json:"areaCode"`
	PhoneNum  string `json:"phoneNum"`
	Sms       string `json:"sms"`
	SendTo    string `json:"sendTo"`
}

// AuthSmsStep 验证短信
func (s *speed) AuthSmsStep(uin string, step string, sig string, rand string, ticket string, areacode string, mobile string, anonymousMobile string) (*SmsVerify, error) {
	arr := strings.Split(s.addr, ":")
	url := "http://" + arr[0] + ":3500/?&type=Android/&pe=authentication" + step + "/"
	data := map[string]string{
		"Uin":             uin,
		"sig":             sig,
		"randstr":         rand,
		"ticket":          ticket,
		"mobile":          mobile,
		"anonymousmobile": anonymousMobile,
		"areacode":        areacode,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return nil, errors.New("发送请求失败")
	}
	var sms *SmsVerify
	err = json.Unmarshal([]byte(res.Body), &sms)
	if err != nil || res == nil {
		return nil, errors.New("解析数据失败")
	}
	return sms, nil
}

const (
	Error       = iota //登陆失败
	Success            //登陆成功
	NeedCaptcha        //需要验证滑动验证码
	NeedCode           //需要验证短信验证码
	PassError          //密码错误
	EnvError           //环境异常
	UinLock            //账号冻结
	NotCreate          //未添加QQ
	NeedSendSms        //需要发送短信
)

type captcha struct {
	HkUrl   string `json:"HkUrl"`
	Stype   string `json:"Stype"`
	Content string `json:"fromContent"`
}

// getQrCode 获取设备锁二维码
func (s *speed) getQrCode(verifyUrl string) string {
	if !strings.Contains(verifyUrl, "uin=") || !strings.Contains(verifyUrl, "&sig=") {
		return ""
	}
	//解析参数
	uin := str.NewStringTool(verifyUrl).GetBetween("uin=", "&sig=")
	sig := str.NewStringTool(verifyUrl).GetBetween("&sig=", "-1")
	//构建url
	url := "https://ti.qq.com/safe/verify/proxy/domain/oidb.tim.qq.com/v3/oidbinterface/oidb_0xc9e_4?uid=" + uin + "&getqrcode=1&sdkappid=39998&actype=2&bkn="
	req, err := http.NewRequest("POST", url, strings.NewReader("{\"str_dev_auth_token\":\""+sig+"\",\"uint32_flag\":1}"))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Cookie", "domainId=338; uin=o"+uin+"; p_uin=o"+uin)
	req.Header.Set("Origin", "https://ti.qq.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("qname-service", "1104449:131072")
	req.Header.Set("qname-space", "Production")
	req.Header.Set("Referer", verifyUrl)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.124 Safari/537.36 Edg/102.0.1245.41")
	resp, err := getClient().Do(req)
	if err != nil {
		return ""
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	//解析返回值
	type result struct {
		Uint32GuaranteeStatus int    `json:"uint32_guarantee_status"`
		StrUrl                string `json:"str_url"`
		ActionStatus          string `json:"ActionStatus"`
		ErrorCode             int    `json:"ErrorCode"`
		ErrorInfo             string `json:"ErrorInfo"`
	}
	var res *result
	err = json.Unmarshal(bytes, &res)
	if err != nil || res == nil {
		return ""
	}
	if res.ErrorCode != 0 {
		return ""
	}
	d := "https://ti.qq.com/safe/scanresult?str_url=" + res.StrUrl + "&envfrom=&verify_id=&verify_scene=&uin=" + uin
	return d
}

// login 登陆账号
func (s *speed) login(uin string, code string) (int, error, string, string) {
	//添加账号时候等待
	server2.GetIns().ServerWaitLock("add_" + s.addr)
	//判断qq是否并发登陆
	if server2.GetIns().ServerIsLock("uin_" + uin) {
		return Error, errors.New("操作频繁,请稍后重试"), "", ""
	} else {
		server2.GetIns().ServerLock("uin_" + uin)
	}
	defer server2.GetIns().ServerUnLock("uin_" + uin) //解除并发锁
	//阻止全局新增
	server2.GetIns().ServerInc("login_"+s.addr, 120)
	defer server2.GetIns().ServerDec("login_"+s.addr, 120)

	//等待计数器完成
	server2.GetIns().ServerNumWaitZero("add_" + s.addr)

	url := "/?&type=Android/&pe=login/"
	data := map[string]string{
		"Uin":   uin,
		"code":  code,
		"place": s.proxy,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return Error, errors.New("服务正在升级中,请稍后重试"), "", ""
	}
	var state *captcha
	err = json.Unmarshal([]byte(res.State), &state)
	if err != nil || state == nil {
		return Error, errors.New("解析服务器结果错误"), "", ""
	}
	if state.Stype == "-2" {
		return NotCreate, errors.New("系统正在大量挂机中,请稍后重试"), "", ""
	}
	if state.Stype == "0" || state.Stype == "-1" || res.State == "0" {
		return Success, nil, "", strconv.FormatInt(res.Uin, 10)
	}
	if state.Stype == "6" || state.Stype == "10" {
		return Error, errors.New("登陆失败"), "", ""
	}
	if state.Stype == "7" && state.HkUrl == "" {
		return EnvError, errors.New("登陆环境异常,请切换网络重试"), "", ""
	}
	if state.Stype == "7" && state.HkUrl != "" {
		return NeedSendSms, errors.New("请验证安全手机"), state.HkUrl, ""
	}
	if state.Stype == "8" {
		return UinLock, errors.New("账号已被冻结"), "", ""
	}
	if state.Stype == "9" {
		return PassError, errors.New("账号密码错误"), "", ""
	}
	if state.Stype == "11" {
		return NeedCaptcha, errors.New("请验证滑动验证码"), state.HkUrl, ""
	}
	if state.Stype == "12" || state.Stype == "13" {
		return NeedCode, errors.New("请提交设备锁验证码"), s.getQrCode(state.HkUrl), ""
	}
	return Error, errors.New("未知状态"), "", ""
}

// checkLogin 判断是否登录
func (s *speed) checkLogin(uin string) bool {
	return uin == ""
}
