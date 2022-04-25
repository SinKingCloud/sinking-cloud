package uin_tool

import (
	"encoding/json"
	"errors"
)

// add 添加账号
func (s *speed) add(uin string, pwd string) bool {
	url := "/?&type=Android/&pe=into/"
	data := map[string]string{
		"Uin":      uin,
		"Pass":     pwd,
		"place":    s.area,
		"brand":    "",
		"model":    "",
		"Protocol": "0",
	}
	res, err := s.sendRequest(url, data)
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
	PassError          //密码错误
	EnvError           //环境异常
	UinLock            //账号冻结
	NotCreate          //未创建QQ
)

type captcha struct {
	HkUrl string `json:"HkUrl"`
	Stype string `json:"Stype"`
}

// login 登陆账号
func (s *speed) login(uin string, code string) (int, error, string) {
	url := "/?&type=Android/&pe=login/"
	data := map[string]string{
		"Uin":  uin,
		"Code": code,
	}
	res, err := s.sendRequest(url, data)
	if err != nil {
		return Error, err, ""
	}
	var state *captcha
	err = json.Unmarshal([]byte(res.State), &state)
	if err != nil || state == nil {
		return Error, errors.New("解析服务器结果错误"), ""
	}
	if state.Stype == "-2" {
		return NotCreate, errors.New("系统不存在此QQ"), ""
	}
	if state.Stype == "0" || state.Stype == "-1" || res.State == "0" {
		return Success, nil, ""
	}
	if state.Stype == "6" || state.Stype == "10" {
		return Error, errors.New("登陆失败"), ""
	}
	if state.Stype == "7" {
		return EnvError, errors.New("登陆环境异常"), ""
	}
	if state.Stype == "8" {
		return UinLock, errors.New("账号已被冻结"), ""
	}
	if state.Stype == "9" {
		return PassError, errors.New("账号密码错误"), ""
	}
	if state.Stype == "12" {
		return NeedCode, errors.New("请提交设备锁验证码"), ""
	}
	if state.Stype == "11" {
		return NeedCaptcha, errors.New("请验证滑动验证码"), state.HkUrl
	}
	return Error, errors.New("未知状态"), ""
}
