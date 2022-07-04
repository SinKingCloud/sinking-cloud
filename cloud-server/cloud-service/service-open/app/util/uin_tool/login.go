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
)

// add 添加账号
func (s *speed) add(uin string, pwd string, phone string, proto int) bool {
	//上锁
	server2.GetIns().ServerLock(s.addr)
	defer server2.GetIns().ServerUnLock(s.addr)
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
		"protocol": p,
		"phone":    phone,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

const (
	Error       = iota //登陆失败
	Success            //登陆成功
	NeedCaptcha        //需要验证滑动验证码
	NeedCode           //需要验证短信验证码
	NeedQrCode         //需要扫码验证
	PassError          //密码错误
	EnvError           //环境异常
	UinLock            //账号冻结
	NotCreate          //未添加QQ
)

type captcha struct {
	HkUrl   string `json:"HkUrl"`
	Stype   string `json:"Stype"`
	Content string `json:"fromContent"`
}

func (s *speed) getQrCode(verifyUrl string) string {
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
	//阻塞锁
	server2.GetIns().ServerWaitLock(s.addr)
	//判断qq是否并发登陆
	if server2.GetIns().ServerIsLock(uin) {
		return Error, errors.New("操作频繁,请稍后重试"), "", ""
	}
	//解除并发锁
	defer server2.GetIns().ServerUnLock(uin)
	url := "/?&type=Android/&pe=login/"
	data := map[string]string{
		"Uin":   uin,
		"code":  code,
		"place": s.proxy,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return Error, errors.New("服务器连接失败,请稍后重试"), "", ""
	}
	var state *captcha
	err = json.Unmarshal([]byte(res.State), &state)
	if err != nil || state == nil {
		return Error, errors.New("解析服务器结果错误"), "", ""
	}
	if state.Stype == "-2" {
		return NotCreate, errors.New("系统不存在此QQ"), "", ""
	}
	if state.Stype == "0" || state.Stype == "-1" || res.State == "0" {
		return Success, nil, "", strconv.FormatInt(res.Uin, 10)
	}
	if state.Stype == "6" || state.Stype == "10" {
		return Error, errors.New("登陆失败"), "", ""
	}
	if state.Stype == "7" && state.HkUrl == "" {
		return EnvError, errors.New("登陆环境异常"), "", ""
	}
	if state.Stype == "7" && state.HkUrl != "" {
		return NeedQrCode, errors.New("请使用QQ扫描二维码验证"), s.getQrCode(state.HkUrl), ""
	}
	if state.Stype == "8" {
		return UinLock, errors.New("账号已被冻结"), "", ""
	}
	if state.Stype == "9" {
		return PassError, errors.New("账号密码错误"), "", ""
	}
	if state.Stype == "12" || state.Stype == "13" {
		return NeedCode, errors.New("请提交设备锁验证码"), s.getQrCode(state.HkUrl), ""
	}
	if state.Stype == "11" {
		return NeedCaptcha, errors.New("请验证滑动验证码"), state.HkUrl, ""
	}
	return Error, errors.New("未知状态"), "", ""
}
