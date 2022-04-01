package uin_tool

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
	url := "/?&type=Android/&pe=limixiuo/"
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
	url := "/?&type=Pc/&pe=weishi/"
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
	url := "/?&type=Pc/&pe=guanjia/"
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
