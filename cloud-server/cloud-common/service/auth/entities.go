package auth

// CheckCaptchaParam 判断滑块验证码
type CheckCaptchaParam struct {
	Key string //验证码标识
	X   int    //滑块x坐标
	Y   int    //滑块y坐标
}

// CheckSmsParam 判断短信验证码
type CheckSmsParam struct {
	Phone string //手机号
	Code  string //验证码
}

// CheckMailParam 判断邮件验证码
type CheckMailParam struct {
	Mail string //邮箱
	Code string //验证码
}

// WebUserLoginParam 网站用户登陆参数
type WebUserLoginParam struct {
	Account   string
	Password  string
	Code      string
	UserAgent string
	Ip        string
}

// SystemUserLoginParam 系统用户登陆参数
type SystemUserLoginParam struct {
	Account   string
	Password  string
	Code      string
	UserAgent string
	Ip        string
}
