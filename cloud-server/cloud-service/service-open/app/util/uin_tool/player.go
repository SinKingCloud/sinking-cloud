package uin_tool

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	url := "/?&type=Pc/&pe=fangke/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	return res.State == "0"
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
func (s *speed) zoneClockSpeed(uin string) bool {
	url := "/?&type=Android/&pe=qzone/"
	data := map[string]string{
		"Uin": uin,
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
func (s *speed) pcLogin(uin string) bool {
	url := "/?&type=pc/&pe=login/"
	data := map[string]string{
		"Uin":   uin,
		"place": s.proxy,
	}
	res, err := s.sendRequest(url, data, -1)
	if err != nil {
		return false
	}
	var state *captcha
	err = json.Unmarshal([]byte(res.State), &state)
	if err != nil || state == nil {
		return false
	}
	return state.Stype == "11"
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
	if res.State == "-2" {
		if s.pcAdd(uin) {
			res, err = s.sendRequest(url, data, -1)
			if err != nil {
				return false
			}
		}
	}
	return res.State == "0"
}

// pcForceSpeed 电脑加速
func (s *speed) pcForceSpeed(uin string) bool {
	r := s.pcLogin(uin)
	if !r {
		r = s.pcLogin(uin)
		if !r {
			r = s.pcLogin(uin)
		}
	}
	if r {
		return s.pcCode(uin)
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
	return s.pcLogin(uin)
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
