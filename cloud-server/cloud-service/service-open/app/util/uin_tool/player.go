package uin_tool

import (
	"encoding/json"
)

// sportSpeed 运动加速
func (s *speed) sportSpeed(uin string) bool {
	url := "/?&type=Android/&pe=yundong/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data)
	if err != nil {
		return false
	}
	return res.State == "0" || res.State == "-10001" || res.State == "-10003"
}

// microShowSpeed 厘米秀加速
func (s *speed) microShowSpeed(uin string) bool {
	url := "/?&type=Android/&pe=Limixiu/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data)
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
	res, err := s.sendRequest(url, data)
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
	res, err := s.sendRequest(url, data)
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
	res, err := s.sendRequest(url, data)
	if err != nil {
		return false
	}
	return res.State != "-1"
}

// zoneSpeed 空间加速
func (s *speed) zoneSpeed(uin string) bool {
	url := "/?&type=Pc/&pe=fangke/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data)
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
	res, err := s.sendRequest(url, data)
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
	res, err := s.sendRequest(url, data)
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
	res, err := s.sendRequest(url, data)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// pcSpeed 电脑加速
func (s *speed) pcSpeed(uin string) bool {
	url := "/?&type=Android/&pe=pc/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// friendSpeed 加好友加速
func (s *speed) friendSpeed(uin string) bool {
	url := "/?&type=Android/&pe=Friends/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data)
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
	res, err := s.sendRequest(url, data)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// pcLower pc下线
func (s *speed) pcLower(uin string) bool {
	url := "/?&type=Pc/&pe=lower/"
	data := map[string]string{
		"Uin": uin,
	}
	res, err := s.sendRequest(url, data)
	if err != nil {
		return false
	}
	return res.State == "0"
}

// getSpeedStatus 获取加速状态
func (s *speed) getSpeedStatus(uin string) (d *SpeedStatus, e error) {
	url := "/?&type=Android/&pe=querygrade/"
	data := map[string]interface{}{
		"Uin":  uin,
		"code": 0,
	}
	res, e := s.sendRequest(url, data)
	if e != nil {
		return
	}
	if res.State == "" {
		return nil, nil
	}
	e = json.Unmarshal([]byte(res.State), &d)
	return
}
