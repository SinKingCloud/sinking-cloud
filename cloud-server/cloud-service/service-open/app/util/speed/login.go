package speed

import "errors"

// add 添加账号
func (s *speed) add(uin string, pwd string) bool {
	url := "/?&type=Android/&pe=into/"
	data := map[string]string{
		"Uin":      uin,
		"Pass":     pwd,
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
)

// login 登陆账号
func (s *speed) login(uin string, code string) (int, error) {
	url := "/?&type=Android/&pe=login/"
	data := map[string]string{
		"Uin":  uin,
		"Code": code,
	}
	res, err := s.sendRequest(url, data)
	if err != nil {
		return Error, err
	}
	if res.State == "-1" {
		return Error, errors.New("系统不存在此QQ")
	}
	if res.State == "0" {
		return Success, nil
	}
	if res.State == "6" || res.State == "10" {
		return Error, errors.New("登陆失败")
	}
	if res.State == "7" {
		return EnvError, errors.New("登陆环境异常")
	}
	if res.State == "8" {
		return UinLock, errors.New("账号已被冻结")
	}
	if res.State == "9" {
		return PassError, errors.New("账号密码错误")
	}
	if res.State == "11" {
		return NeedCaptcha, errors.New("请验证滑动验证码")
	}
	if res.State == "12" {
		return NeedCode, errors.New("请提交设备锁验证码")
	}
	return Error, errors.New("未知状态")
}
