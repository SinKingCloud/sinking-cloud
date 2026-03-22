package service

type AuthService interface {
	GetCaptcha(key string) (map[string]interface{}, error)                                                         //获取滑块验证码
	CheckCaptcha(key string, x int, y int) bool                                                                    //判断滑块验证码
	CheckSms(phone string, code string) bool                                                                       //判断短信验证码
	SendSms(phone string) bool                                                                                     //发送短信验证码
	CheckMail(mail string, code string) bool                                                                       //判断邮件验证码
	SendMail(mail string) bool                                                                                     //发送邮件验证码
	Login(account string, isAdmin bool, password string, code string, userAgent string, ip string) (string, error) //账号登录或注册
}
