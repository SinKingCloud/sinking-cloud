package uin_tool

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// sign 构造参数
type sign struct {
	uin        string //QQ账号
	sKey       string //sKey
	psKey      string //psKey
	proxy      string //代理地址
	cookie     string //cookie
	gtk        string //gtk
	gtk2       string //gtk2
	httpClient *http.Client
}

// NewSign 实例化sign
func NewSign(uin string, sKey string, psKey string) *sign {
	s := &sign{
		uin:   uin,
		sKey:  sKey,
		psKey: psKey,
	}
	s.gtk2 = s.getGTK2(sKey)
	s.gtk = strconv.FormatInt(s.getGTK(psKey), 10)
	s.cookie = "pt2gguin=o0" + uin + "; uin=o0" + uin + "; skey=" + sKey + "; p_skey=" + psKey + "; p_uin=o0" + uin
	return s
}

// SetProxy 设置代理
func (s *sign) SetProxy(proxy string) *sign {
	s.proxy = proxy
	return s
}

// getClient 获取请求client
func (s *sign) getClient() *http.Client {
	if s.httpClient == nil {
		s.httpClient = &http.Client{
			Timeout: 60 * time.Second,
		}
		if s.proxy != "" {
			proxy, err := url.Parse(s.proxy)
			if err == nil {
				s.httpClient.Transport = &http.Transport{
					Proxy: http.ProxyURL(proxy),
				}
			}
		}
	}
	return s.httpClient
}

// signResponse 响应对象
type signResponse struct {
	Body   string
	Header http.Header
}

// request 发送request请求
func (s *sign) request(url string, method string, headers map[string]string, data string, isGzip bool) (*signResponse, error) {
	var reader io.Reader
	if data != "" {
		reader = strings.NewReader(data)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := s.getClient().Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	var body []byte
	if isGzip {
		r, e := gzip.NewReader(resp.Body)
		if e != nil {
			return nil, e
		}
		b, e := ioutil.ReadAll(r)
		if e != nil {
			return nil, e
		}
		body = b
	} else {
		b, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return nil, e
		}
		body = b
	}
	res := &signResponse{
		Body:   string(body),
		Header: resp.Header.Clone(),
	}
	return res, nil
}

// Bin2hex 转16进制字节
func (*sign) Bin2hex(str string) string {
	return hex.EncodeToString([]byte(str))
}

// Ord ord()
func (*sign) Ord(char string) int {
	r, _ := utf8.DecodeRune([]byte(char))
	return int(r)
}

// HexDec HexDec()
func (*sign) HexDec(str string) int64 {
	num, e := strconv.ParseInt(str, 16, 0)
	if e != nil {
		return 0
	}
	return num
}

func (s *sign) Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *sign) getGTK(sKey string) int64 {
	hash := int64(5381)
	for i := 0; i < len(sKey); i++ {
		hash += (hash << 5 & 2147483647) + int64(s.Ord(string(sKey[i])))&2147483647
		hash &= 2147483647
	}
	return hash & 2147483647
}

func (s *sign) getGTK2(sKey string) string {
	salt := int64(5381)
	md5key := "tencentQQVIP123443safde&!%^%1282"
	var hash []int64
	hash = append(hash, salt<<5)
	for i := 0; i < len(sKey); i++ {
		str := string(sKey[i])
		ASCIICode := s.HexDec(s.Bin2hex(str))
		hash = append(hash, (salt<<5)+ASCIICode)
		salt = ASCIICode
	}
	str := ""
	for _, v := range hash {
		str += strconv.FormatInt(v, 10)
	}
	return s.Md5(str + md5key)
}

// VipSign 会员签到
func (s *sign) VipSign() {
	h := map[string]string{
		"user-agent": "Mozilla/5.0 (Linux; U; Android 4.4.1; zh-cn) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.5 Mobile Safari/533.1",
		"referer":    "https://vip.qq.com/",
		"cookie":     s.cookie,
	}
	u := "https://iyouxi3.vip.qq.com/ams3.0.php?_c=page&actid=79968&format=json&g_tk=" + s.gtk2 + "&cachetime=" + strconv.FormatInt(time.Now().Unix(), 10)
	_, _ = s.request(u, "GET", h, "", false) //会员面板签到
	u = "https://iyouxi3.vip.qq.com/ams3.0.php?_c=page&actid=403490&rand=0.27489888" + strconv.FormatInt(time.Now().Unix(), 10) + "&g_tk=" + s.gtk2 + "&format=json"
	_, _ = s.request(u, "GET", h, "", false) //会员电脑端签到
	u = "https://iyouxi3.vip.qq.com/ams3.0.php?_c=page&actid=23074&format=json&g_tk=" + s.gtk2
	_, _ = s.request(u, "GET", h, "", false) //会员积分手机端签到
	u = "https://iyouxi3.vip.qq.com/ams3.0.php?_c=page&actid=288341&format=json&g_tk=" + s.gtk2
	_, _ = s.request(u, "GET", h, "", false) //分享积分活动签到
	u = "https://iyouxi4.vip.qq.com/ams3.0.php?pvsrc=&s_p=0|http|&s_v=0&ozid=&vipid=&actid=200084&g_tk=" + s.gtk2
	_, _ = s.request(u, "GET", h, "", false) //超级会员每日10萌点领取
}

// BigVipSign 大会员签到
func (s *sign) BigVipSign() {
	h := map[string]string{
		"user-agent": "Mozilla/5.0 (Linux; U; Android 4.4.1; zh-cn) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.5 Mobile Safari/533.1",
		"cookie":     s.cookie,
	}
	u := "https://vip.qzone.qq.com/fcg-bin/v2/fcg_vip_task_checkin?t=0." + strconv.FormatInt(time.Now().Unix(), 10) + "082161&g_tk=" + s.gtk
	post := "appid=qq_big_vip&op=CheckIn&uin=" + s.uin + "&format=json&inCharset=utf-8&outCharset=utf-8"
	_, _ = s.request(u, "POST", h, post, false) //大会员签到

	u = "https://h5.qzone.qq.com/qzone/visitor?_wv=3&_wwv=1024&_proxy=1"
	_, _ = s.request(u, "GET", h, "", false) //登陆访客
	u = "https://h5.qzone.qq.com/webapp/json/QQBigVipTask/CompleteTask?t=0." + strconv.FormatInt(time.Now().Unix(), 10) + "906319&g_tk=" + s.gtk
	post = "outCharset=utf-8&iAppId=0&llTime=" + strconv.FormatInt(time.Now().Unix(), 10) + "&format=json&iActionType=6&strUid=" + s.uin + "&uin=" + s.uin + "&inCharset=utf-8"
	_, _ = s.request(u, "POST", h, post, false) //登陆访客

	u = "https://h5.qzone.qq.com/v2/vip/dhy/trpc/home/ReportShare"
	h["Content-Type"] = "application/json; charset=utf-8"
	post = "{\"g_tk\":" + s.gtk + ",\"head\":{\"uid\":\"" + s.uin + "\",\"skey\":\"" + s.sKey + "\"},\"uid\":\"" + s.uin + "\"}"
	_, _ = s.request(u, "POST", h, post, false) //大会员分享

	u = "https://h5.qzone.qq.com/v2/vip/dhy/trpc/home/ReportRelationship"
	_, _ = s.request(u, "POST", h, post, false) //大会员师徒打卡
}

// YellowDiaSin 黄钻签到
func (s *sign) YellowDiaSin() {
	h := map[string]string{
		"user-agent": "Mozilla/5.0 (Linux; U; Android 4.4.1; zh-cn) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.5 Mobile Safari/533.1",
		"cookie":     s.cookie,
		"referer":    "https://h5.qzone.qq.com/vipinfo/index?plg_nld=1&source=qqmail&plg_auth=1&plg_uin=1&_wv=3&plg_dev=1&plg_nld=1&aid=jh&_bid=368&plg_usr=1&plg_vkey=1&pt_qzone_sig=1",
	}
	u := "https://h5.qzone.qq.com/proxy/domain/vip.qzone.qq.com/fcg-bin/v2/fcg_mobile_vip_site_checkin?t=0.89457" + strconv.FormatInt(time.Now().Unix(), 10) + "&g_tk=" + s.gtk + "&qzonetoken=__vip_score"
	_, _ = s.request(u, "POST", h, "uin="+s.uin+"&format=json", false) //黄钻签到
}
