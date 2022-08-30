package uin_tool

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
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
	superKey   string //superKey
	proxy      string //代理地址
	cookie     string //cookie
	gtk        string //gtk
	gtk2       string //gtk2
	httpClient *http.Client
}

// NewSignTool 实例化sign
func NewSignTool(uin string, proxy string, sKey string, psKey string, superKey string) *sign {
	s := &sign{
		uin:      uin,
		sKey:     sKey,
		psKey:    psKey,
		superKey: superKey,
		proxy:    proxy,
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
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		tran := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		if s.proxy != "" {
			proxy, err := url.Parse("http://" + s.proxy)
			if err == nil {
				tran.Proxy = http.ProxyURL(proxy)
			}
		}
		s.httpClient.Transport = tran
	}
	return s.httpClient
}

// signResponse 响应对象
type signResponse struct {
	Body   string
	Header http.Header
}

// byteToString byte转string
func (s *sign) byteToString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

// getRandom 获取随机数
func (s *sign) getRandom(l int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	strList := []byte(str)

	result := []byte{}
	i := 0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i < l {
		n := strList[r.Intn(len(strList))]
		result = append(result, n)
		i = i + 1
	}
	return string(result)
}

func (s *sign) GzipEncode(input []byte) ([]byte, error) {
	// 创建一个新的 byte 输出流
	var buf bytes.Buffer
	// 创建一个新的 gzip 输出流
	gzipWriter := gzip.NewWriter(&buf)
	// 将 input byte 数组写入到此输出流中
	_, err := gzipWriter.Write(input)
	if err != nil {
		_ = gzipWriter.Close()
		return nil, err
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}
	// 返回压缩后的 bytes 数组
	return buf.Bytes(), nil
}

func (s *sign) GzipDecode(input []byte) ([]byte, error) {
	// 创建一个新的 gzip.Reader
	bytesReader := bytes.NewReader(input)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return nil, err
	}
	defer func() {
		// defer 中关闭 gzipReader
		_ = gzipReader.Close()
	}()
	buf := new(bytes.Buffer)
	// 从 Reader 中读取出数据
	if _, err := buf.ReadFrom(gzipReader); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
	var body string
	if isGzip {
		r, e := gzip.NewReader(resp.Body)
		if e != nil {
			return nil, e
		}
		b, e := ioutil.ReadAll(r)
		if e != nil {
			return nil, e
		}
		body = string(b)
	} else {
		b, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return nil, e
		}
		body = string(b)
	}
	res := &signResponse{
		Body:   body,
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

func (s *sign) getToken(superKey string) int64 {
	if superKey == "" {
		return 0
	}
	var hash int64
	for i := 0; i < len(superKey); i++ {
		n := int64(s.Ord(string(superKey[i])))
		hash = (hash*33 + n) % int64(4294967296)
	}
	return hash
}

func (s *sign) getToken2(superKey string) int64 {
	if superKey == "" {
		return 0
	}
	var hash int64
	for i := 0; i < len(superKey); i++ {
		n := int64(s.Ord(string(superKey[i])))
		hash += (hash << 5 & 2147483647) + n&2147483647
		hash &= 2147483647
	}
	return hash & 2147483647
}

// cutStr 获取文本中间字符
func (s *sign) cutStr(str string, start string, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[len(start):m])
	return str
}

// getCookie 获取cookie
func (s *sign) getCookie(cookies []string) map[string]string {
	data := make(map[string]string)
	for _, cookie := range cookies {
		arr := strings.Split(cookie, ";")
		for _, v := range arr {
			a := strings.Split(v, "=")
			k := strings.ToLower(a[0])
			if len(a) >= 2 && a[1] != "" && k != "path" && k != "domain" && k != "expires" {
				data[a[0]] = a[1]
			}
		}
	}
	return data
}

// getCookieString 获取cookie
func (s *sign) getCookieString(cookies map[string]string) string {
	cookie := ""
	for k, v := range cookies {
		cookie += k + "=" + v + "; "
	}
	return cookie
}

// getPtSigX 授权获取key
func (s *sign) getPtSigX(daid string, appid string) string {
	if s.superKey == "" {
		return ""
	}
	token := s.getToken(s.superKey)
	u := "https://ssl.ptlogin2.qq.com/pt4_auth?daid=" + daid + "&appid=" + appid + "&auth_token=" + strconv.FormatInt(s.getToken(strconv.FormatInt(token, 10)), 10)
	r := "https://ui.ptlogin2.qq.com/cgi-bin/login"
	c := "superuin=o0" + s.uin + "; superkey=" + s.superKey + "; supertoken=" + strconv.FormatInt(token, 10) + ";"
	h := map[string]string{
		"user-agent": "Mozilla/5.0 (Linux; U; Android 4.4.1; zh-cn) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.5 Mobile Safari/533.1",
		"referer":    r,
		"cookie":     c,
	}
	res, e := s.request(u, "GET", h, "", false)
	if e != nil {
		return ""
	}
	return s.cutStr(res.Body, "ptsigx=", "&")
}

// VipSign 会员签到
func (s *sign) VipSign() {
	h := map[string]string{
		"user-agent": "Mozilla/5.0 (Linux; U; Android 4.4.1; zh-cn) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.5 Mobile Safari/533.1",
		"referer":    "https://vip.qq.com/",
		"cookie":     s.cookie,
		"origin":     "https://iyouxi3.vip.qq.com/",
	}
	u := "https://iyouxi3.vip.qq.com/ams3.0.php?_c=page&actid=79968&format=json&g_tk=" + s.gtk2 + "&cachetime=" + strconv.FormatInt(time.Now().Unix(), 10)
	_, _ = s.request(u, "GET", h, "", false) //会员面板签到
	u = "https://iyouxi3.vip.qq.com/ams3.0.php?_c=page&actid=403490&rand=0.27489888" + strconv.FormatInt(time.Now().Unix(), 10) + "&g_tk=" + s.gtk2 + "&format=json"
	_, _ = s.request(u, "GET", h, "", false) //会员电脑端签到
	u = "https://iyouxi3.vip.qq.com/ams3.0.php?_c=page&actid=23074&format=json&g_tk=" + s.gtk2
	_, _ = s.request(u, "GET", h, "", false) //会员积分手机端签到
	u = "https://iyouxi3.vip.qq.com/ams3.0.php?_c=page&actid=288341&format=json&g_tk=" + s.gtk2
	_, _ = s.request(u, "GET", h, "", false) //分享积分活动签到
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

// BlueDiaSin 蓝钻签到
func (s *sign) BlueDiaSin() {
	sig := s.getPtSigX("176", "21000110")
	if sig != "" {
		u := "https://ptlogin2.gamevip.qq.com/check_sig?uin=" + s.uin + "&ptsigx=" + sig + "&daid=176&pt_login_type=4&service=pt4_auth&pttype=2&regmaster=&aid=21000110&s_url=http%3A%2F%2Fgamevip.qq.com%2F"
		h := map[string]string{
			"user-agent": "Mozilla/5.0 (Linux; U; Android 4.4.1; zh-cn) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.5 Mobile Safari/533.1",
			"referer":    "https://ui.ptlogin2.qq.com/cgi-bin/login",
		}
		res, e := s.request(u, "GET", h, "", false)
		if e == nil {
			u = "https://app.gamevip.qq.com/cgi-bin/gamevip_sign/GameVip_SignIn?format=json&g_tk=" + s.gtk + "&_=" + strconv.FormatInt(time.Now().Unix(), 10) + "034"
			h = map[string]string{
				"user-agent": "Mozilla/5.0 (Linux; U; Android 4.4.1; zh-cn) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.5 Mobile Safari/533.1",
				"referer":    "https://gamevip.qq.com/sign_pop/sign_pop_v2.html",
				"cookie":     s.getCookieString(s.getCookie(res.Header.Values("Set-Cookie"))) + "DomainID=176;",
			}
			_, _ = s.request(u, "GET", h, "", false) //蓝钻签到

			u = "https://app.gamevip.qq.com/cgi-bin/gamevip_sign/GameVip_Lottery?format=json&g_tk=" + s.gtk + "&_=" + strconv.FormatInt(time.Now().Unix(), 10) + "034"
			h["referer"] = u
			_, _ = s.request(u, "GET", h, "", false) //蓝钻抽奖
		}
	}

}

// DaleDouSign 大乐斗签到
func (s *sign) DaleDouSign() {
	u := "https://fight.pet.qq.com/cgi-bin/petpk?cmd=award&op=1&type=0"
	h := map[string]string{
		"user-agent": "Mozilla/5.0 (Linux; U; Android 4.4.1; zh-cn) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.5 Mobile Safari/533.1",
		"cookie":     s.cookie,
		"referer":    u,
	}
	_, _ = s.request(u, "GET", h, "", false)
}

// QQGamSign QQ游戏签到
func (s *sign) QQGamSign() {
	u := "https://social.minigame.qq.com/cgi-bin/social/welcome_panel_operate?format=json&cmd=2&uin=" + s.uin + "&g_tk=" + s.gtk
	h := map[string]string{
		"Accept-Encoding": "gzip,deflate,sdch",
		"user-agent":      "Mozilla/5.0 (Linux; U; Android 4.4.1; zh-cn) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.5 Mobile Safari/533.1",
		"cookie":          s.cookie,
		"referer":         "https://minigame.qq.com/appdir/social/cloudHall/src/index/welcome.html",
	}
	_, _ = s.request(u, "GET", h, "", true)
	u = "https://social.minigame.qq.com/cgi-bin/social/welcome_panel_operate?format=json&cmd=1&uin=" + s.uin + "&g_tk=" + s.gtk
	_, _ = s.request(u, "GET", h, "", true)
	u = "https://social.minigame.qq.com/cgi-bin/social/CheckInPanel_Operate?Cmd=CheckIn_Operate&g_tk=" + s.gtk
	_, _ = s.request(u, "GET", h, "", true)
}

// GameCenterSign 游戏中心签到
func (s *sign) GameCenterSign() {
	u := "https://w.gamecenter.qq.com/v1/cgi-bin/home-spa/task/receive-task?t=" + strconv.FormatInt(time.Now().Unix(), 10) + "449&g_tk=" + s.gtk + "&p_tk="
	h := map[string]string{
		"user-agent": "Mozilla/5.0 (Linux; Android 10; MI 9 Build/QKQ1.190825.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045330 Mobile Safari/537.36 V1_AND_SQ_8.4.1_1442_YYB_D QQ/8.4.1.4680 NetType/4G WebP/0.3.0 Pixel/1080 StatusBarHeight/76 SimpleUISwitch/1 QQTheme/2921",
		"cookie":     s.cookie,
		"referer":    "https://speed.gamecenter.qq.com/pushgame/v1/home/task?_wv=18950115&_wwv=393",
	}
	_, _ = s.request(u, "POST", h, "taskid=1&taskType=1&taskSubType=0&appid=&area=2&partition=0&roleId=", false)
	u = "https://info.gamecenter.qq.com/cgi-bin/gc_my_tab_async_fcgi?merge=1&ver=0&st=" + strconv.FormatInt(time.Now().Unix(), 10) + "746&sid=&uin=" + s.uin + "&number=0&path=489&plat=qq&gamecenter=1&_wv=1031&_proxy=1&gc_version=2&ADTAG=gamecenter&notShowPub=1&param=%7B%220%22%3A%7B%22param%22%3A%7B%22platform%22%3A1%2C%22tt%22%3A1%7D%2C%22module%22%3A%22gc_my_tab%22%2C%22method%22%3A%22sign_in%22%7D%7D&g_tk=" + s.gtk
	_, _ = s.request(u, "GET", h, "", false)
}

// WeiYunSign 微云签到
func (s *sign) WeiYunSign() {
	u := "https://h5.weiyun.com/login?go_url=https%3A%2F%2Fh5.weiyun.com%2Fsign_in"
	h := map[string]string{
		"user-agent": "Mozilla/5.0 (Linux; Android 10; MI 9 Build/QKQ1.190825.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045330 Mobile Safari/537.36 V1_AND_SQ_8.4.1_1442_YYB_D QQ/8.4.1.4680 NetType/4G WebP/0.3.0 Pixel/1080 StatusBarHeight/76 SimpleUISwitch/1 QQTheme/2921",
		"cookie":     s.cookie,
		"referer":    "https://h5.weiyun.com/sign_in",
	}
	_, _ = s.request(u, "GET", h, "", false)
}

// QQMusicSign QQ音乐签到
func (s *sign) QQMusicSign() {
	u := "https://u.y.qq.com/cgi-bin/musicu.fcg"
	h := map[string]string{
		"content-type": "application/json",
		"user-agent":   "Mozilla/5.0 (Linux; Android 10; MI 9 Build/QKQ1.190825.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045330 Mobile Safari/537.36 V1_AND_SQ_8.4.1_1442_YYB_D QQ/8.4.1.4680 NetType/4G WebP/0.3.0 Pixel/1080 StatusBarHeight/76 SimpleUISwitch/1 QQTheme/2921",
		"cookie":       s.cookie,
		"referer":      "https://u.y.qq.com/cgi-bin/musicu.fcg",
	}

	post := "{\"req_0\":{\"module\":\"UserGrow.UserGrowScore\",\"method\":\"receive_score\",\"param\":{\"musicid\":\"" + s.uin + "\",\"type\":15}},\"comm\":{\"g_tk\":" + s.gtk + ",\"uin\":" + s.uin + ",\"format\":\"json\",\"ct\":23,\"cv\":0}}"
	_, _ = s.request(u, "POST", h, post, false)

	post = "{\"comm\":{\"g_tk\":" + s.gtk + ",\"uin\":" + s.uin + ",\"format\":\"json\",\"inCharset\":\"utf-8\",\"outCharset\":\"utf-8\",\"notice\":0,\"platform\":\"h5\",\"needNewCode\":1,\"uid\":\"\"},\"req_0\":{\"module\":\"music.actCenter.ActCenterSignNewSvr\",\"method\":\"DoSignIn\",\"param\":{\"ActID\":\"PR-Config20200828-31525466015\"}}}"
	_, _ = s.request(u, "POST", h, post, false)

	post = "{\"req_0\":{\"module\":\"UserGrow.UserGrowScore\",\"method\":\"receive_score\",\"param\":{\"musicid\":\"" + s.uin + "\",\"type\":1}},\"comm\":{\"g_tk\":" + s.gtk + ",\"uin\":" + s.uin + ",\"format\":\"json\",\"ct\":23,\"cv\":0}}"
	_, _ = s.request(u, "POST", h, post, false)

	post = "{\"req_0\":{\"module\":\"Radio.RadioLucky\",\"method\":\"clockIn\",\"param\":{\"platform\":2}},\"comm\":{\"g_tk\":" + s.gtk + ",\"uin\":" + s.uin + ",\"format\":\"json\"}}"
	_, _ = s.request(u, "POST", h, post, false)

	u = "https://stat.y.qq.com/android/fcgi-bin/imusic_tj"
	ua := "QQMusic 9020807(android 9)"
	timestamp := time.Now().Unix()
	times := "3599"
	timekey := strings.ToUpper(s.Md5(strconv.FormatInt(timestamp, 10) + times + s.uin + "gk2$Lh-&l4#!4iow"))
	h = map[string]string{
		"user-agent": ua,
	}
	xml := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root><rom>XiaoMi/MIUI/V11</rom><ct>11</ct><cv>9020807</cv><v>9020807</v><mcc>460</mcc><uid>" + strconv.FormatInt(timestamp, 10) + "</uid><mnc>00</mnc><chid>72264</chid><os_ver>9</os_ver><phonetype>MI 1000</phonetype><devicelevel>40</devicelevel><nettype>1023</nettype><qq>" + s.uin + "</qq><cid>228</cid><item cmd=\"1\" optime=\"" + strconv.FormatInt(timestamp, 10) + "\" nettype=\"1023\" QQ=\"" + s.uin + "\" uid=\"" + strconv.FormatInt(timestamp, 10) + "\" os=\"9\" model=\"MI 8\" mcc=\"460\" version=\"9.2.8.7\" songtype=\"1\" playtype=\"4\" from=\"9,30,321,\" dts=\"0\" openstore=\"0\" rec_reason=\"\" crytype=\"5\" paytype=\"3\" abt=\"220_220002\" tjreport=\"\" desktoplyric=\"0\" playdevice=\"0\" outdev=\"0\" url=\"5\" playmode=\"1\" cdn=\"\" cdnip=\"\" hasFirstBuffer=\"3\" filetype=\"3\" err=\"0\"  issoftdecode=\"2\" wait_time=\"292\" player_retry=\"0\" timekey=\"" + timekey + "\" vkey=\"" + s.getRandom(112) + "\" time=\"" + times + "\" errcode=\"\" vip_level=\"0\" songid=\"231769209\" /></root>"
	encode, err := s.GzipEncode([]byte(xml))
	if err == nil {
		_, _ = s.request(u, "GET", h, string(encode), false)
	}
}

func (s *sign) VideoSign() {
	//签到地址
	u := "https://access.video.qq.com/user/auth_login?type=qq&login_flag=1&g_tk=" + s.gtk + "&raw=1&vappid=11059694&vsecret=fdf61a6be0aad57132bc5cdf78ac30145b6cd2c1470b0cfe"
	h := map[string]string{
		"content-type": "application/json",
		"user-agent":   "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/104.0.5112.81",
		"cookie":       s.cookie,
		"referer":      "https://m.v.qq.com/",
	}
	fmt.Println(s.gtk2)
	r, e := s.request(u, "GET", h, "", false)
	fmt.Println(r, e)
}

// VisitLevelPage 查看等级页面
func (s *sign) VisitLevelPage() error {
	req, err := http.NewRequest("GET", "https://ti.qq.com/qqlevel/index?_wv=3&_wwv=1&tab=6&source=15", nil)
	if err != nil {
		return errors.New("发送请求失败")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/104.0.5112.102")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "uin=o"+s.uin+"; skey="+s.sKey+"; p_uin=o"+s.uin+"; p_skey="+s.psKey)
	_, err = s.getClient().Do(req)
	if err != nil {
		return errors.New("发送请求失败")
	}
	return nil
}
