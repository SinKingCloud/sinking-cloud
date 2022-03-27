package constant

const (
	CacheNameWithSmsPhone = "SmsPhone_" //手机验证码储存名称
	CacheTimeWithSmsPhone = 600 * 1000  //手机验证码储存时间

	CacheNameWithExpire = "SmsPhoneExpire_" //手机验证码发送间隔
	CacheTimeWithExpire = 60 * 1000         //发送间隔
)
