package uin_tool

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// sportSpeed 运动加速
func (s *speed) sportSpeed(uin string) bool {
	url := "/?&type=Android/&pe=yundong/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0" || res.State == "-10003"
}

// microShowSpeed 厘米秀加速
func (s *speed) microShowSpeed(uin string) bool {
	url := "/?&type=Android/&pe=Limixiu/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// gameSpeed 游戏加速
func (s *speed) gameSpeed(uin string) bool {
	url := "/?&type=Android/&pe=game/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// wsSpeed 微视加速
func (s *speed) wsSpeed(uin string) bool {
	url := "/?&type=Android/&pe=weishi/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// safeSpeed 管家加速
func (s *speed) safeSpeed(uin string) bool {
	url := "/?&type=Android/&pe=guanjia/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State != "-1"
}

// checkZoneStatus  判断空间状态
func (s *speed) checkZoneStatus(uin string) bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	req, err := http.NewRequest("GET", "https://h5.qzone.qq.com/proxy/domain/r.qzone.qq.com/cgi-bin/qzone_dynamic_v7.cgi?param=848&uin="+uin, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return !strings.Contains(string(bodyText), "-5008")
}

// zoneSpeed 空间加速
func (s *speed) zoneSpeed(uin string) bool {
	//判断空间权限
	if !s.checkZoneStatus(uin) {
		return false
	}
	url := "http://qhzb.xyz:81/api/Outinterface/dlxd?dlid=3294519810&Apikey=ZdTeYUxyD5p8stvU0ZqJXgvkCCbHcQn8&qtype=1&mode=2&sid=&qqnum=" + uin + "&need=10"
	resp, err := getClientWithTimout(30).Get(url)
	if err != nil {
		return false
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		fmt.Println("访客提交:", string(bytes))
	}
	return true
}

// medalSpeed 勋章加速
func (s *speed) medalSpeed(uin string) bool {
	url := "/?&type=Pc/&pe=medal/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// zoneClockSpeed 空间打卡加速
func (s *speed) zoneClockSpeed(uin string, code string) bool {
	url := "/?&type=Android/&pe=qzone/"
	data := map[string]string{
		"Uin":  uin,
		"code": code,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// mobileClockSpeed 手机打卡加速
func (s *speed) mobileClockSpeed(uin string) bool {
	url := "/?&type=Android/&pe=clockin/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// pcOffline 电脑下线
func (s *speed) pcOffline(uin string) bool {
	url := "/?&type=Android/&pe=Offlinepc/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// pcSpeed 电脑加速
func (s *speed) pcLogin(uin string, useProxy bool) (bool, bool, bool) {
	url := "/?&type=pc/&pe=login/"
	data := map[string]string{
		"Uin": uin,
	}
	if useProxy {
		data["place"] = s.proxy
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false, false, false
	}
	var state *captcha
	err = json.Unmarshal([]byte(res.State), &state)
	if err != nil || state == nil {
		return false, false, false
	}
	return state.Stype != "404", state.Stype == "11" || state.Stype == "8", state.Stype == "-1"
}

// pcCode 电脑加速
func (s *speed) pcCode(uin string) bool {
	url := "/?&type=Android/&pe=codepc/"
	data := map[string]string{
		"Uin":   uin,
		"place": s.proxy,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// pcForceSpeed 电脑加速
func (s *speed) pcForceSpeed(uin string) bool {
	useProxy := true
	exists, r, isLogin := s.pcLogin(uin, useProxy)
	if isLogin {
		return true
	}
	if !exists {
		useProxy = false
	}
	if !r {
		exists, r, isLogin = s.pcLogin(uin, useProxy)
		if isLogin {
			return true
		}
		if !r {
			exists, r, isLogin = s.pcLogin(uin, useProxy)
			if isLogin {
				return true
			}
			if !r {
				exists, r, isLogin = s.pcLogin(uin, useProxy)
				if isLogin {
					return true
				}
			}
		}
	}
	if r {
		if s.pcCode(uin) {
			_, b, isLogin := s.pcLogin(uin, useProxy)
			if isLogin {
				return true
			}
			return b
		}
	}
	return false
}

// friendSpeed 加好友加速
func (s *speed) friendSpeed(uin string) bool {
	url := "/?&type=Android/&pe=Friends/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// wordSpeed 小世界加速
func (s *speed) wordSpeed(uin string) bool {
	url := "/?&type=Android/&pe=smallworld/"
	pics := []string{
		"http://worldtj.photo.store.qq.com/psc?/world/O0cFsaTfOlqjIAnYAvw8WkoTvari4*imstS0pam.nlFO3nT*a5oWOFC6I0VYe6knchD5ofMjYyO4Venpw4UXBolsEZluBcB8iaf1m0YEK4k!/b&bo=sASwBLAEsAQRHyg!&wm_text=QOeOi.iAgeWQiSA!",
		"http://worldtj.photo.store.qq.com/psc?/world/O0cFsaTfOlqjIAnYAvw8WhOEdWLO4uOmOY3e9pvVhJd7TTTealNJpTehedk8dthSb.ghI2fcS61E7VNbvKrlPe9.FObzO2Qfk4iSm4BogmE!/b&bo=OAQ4BDgEOAQRHyg!&wm_text=QOeOi.iAgeWQiSA!",
		"http://worldtj.photo.store.qq.com/psc?/world/O0cFsaTfOlqjIAnYAvw8WkoTvari4*imstS0pam.nlFiCyVnMuGDw5ljgTNla6ejjmfRP3f0H8ZNPOjXfHGhtBwlLuoUa1M9gYMH38YEpzY!/b&bo=OAQ4BDgEOAQRHyg!&wm_text=QOeOi.iAgeWQiSA!",
	}
	data := map[string]string{
		"Uin":     uin,
		"picture": pics[rand.Intn(len(pics))],
		"code":    s.RandNum(6),
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// mobileLower 安卓下线
func (s *speed) mobileLower(uin string) bool {
	url := "/?&type=Android/&pe=lower/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// mobileOnline 安卓上线
func (s *speed) mobileOnline(uin string) bool {
	url := "/?&type=Android/&pe=online/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// pcLower pc下线
func (s *speed) pcLower(uin string) bool {
	_, b, _ := s.pcLogin(uin, false)
	return b
}

// getSpeedStatus 获取加速状态
func (s *speed) getSpeedStatus(uin string) (d *SpeedStatus, e error) {
	url := "/?&type=Android/&pe=querygrade/"
	data := map[string]interface{}{
		"Uin":  uin,
		"code": 0,
	}
	res, e := s.sendRequest(url, data, -1)
	if e != nil {
		return
	}
	if res.State == "" {
		return nil, nil
	}
	var t *speedStatus
	e = json.Unmarshal([]byte(res.State), &t)
	if e == nil && t != nil && t.Data != nil {
		d = t.Data
	}
	return
}

// Cookies 获取cookies
func (s *speed) Cookies(uin string, domain string) (string, string) {
	url := "/?&type=Android/&pe=value/"
	data := map[string]interface{}{
		"Uin":  uin,
		"code": domain,
	}
	res, e := s.sendRequest(url, data, -1)
	if e != nil {
		return "", ""
	}
	return res.Skey, res.Pskey
}

// CheckHealth 健康检查
func (s *speed) CheckHealth() bool {
	url := "/?&type=Android/&pe=time/"
	res, e := s.sendRequest(url, "", 2)
	if e != nil {
		return false
	}
	return res.Body != ""
}

// GetCpu cpu使用
func (s *speed) GetCpu() int {
	arr := strings.Split(s.addr, ":")
	url := "http://" + arr[0] + ":1954"
	resp, err := http.Get(url)
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	if err != nil {
		return -1
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1
	}
	num, err := strconv.Atoi(strings.Trim(strings.Replace(string(bytes), "\n", "", -1), ""))
	if err != nil {
		return -1
	}
	return num
}
