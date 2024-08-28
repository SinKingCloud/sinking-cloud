package constant

const (
	AppUinLoginQrCode       = "uin.login.qrcode"       //是否开启扫码更新(0:开启/1:关闭)
	AppUinLoginQrCodeProto  = "uin.login.qrcode.proto" //扫码更新协议(0:IPad/1:AndroidPhone/2:AndroidPad/3:AndroidWatch/4:MacOS)
	AppUinLoginPwd          = "uin.login.pwd"          //是否开启密码更新(0:开启/1:关闭)
	AppUinLoginPwdProto     = "uin.login.pwd.proto"    //密码更新协议(0:IPad/1:AndroidPhone/2:AndroidPad/3:AndroidWatch/4:MacOS)
	UserUinLevel            = "uin.level"              //用户云端挂机折扣等级
	UserUinLevelExpire      = "uin.level.expire"       //用户云端挂机折扣等级到期时间
	UserAuthQrLogin         = "uin.auth.qr_login"      //用户云端挂机是否开通扫码授权权限
	UserAuthDeviceLock      = "uin.auth.device_lock"   //用户云端挂机是否开通设备锁授权权限
	UserAuthCookies         = "uin.auth.cookies"       //用户云端挂机是否开通获取cookies权限
	AppUinSpeedTest         = "speed.test"             //系统云端挂机允许扫码测试
	AppUinSpeedQps          = "speed.qps"              //系统云端挂机扫码测试qps限制
	AppUinSpeedNum          = "speed.num"              //系统云端挂机扫码测试次数限制
	AppUinExperienceOpen    = "experience.open"        //系统云端挂机体验活动开关
	AppUinExperienceDays    = "experience.days"        //系统云端挂机体验活动赠送天数
	AppUinFirstBuyGift      = "first.buy.gift"         //系统云端挂机首次购买赠送天数
	AppUinBuyActivityOpen   = "buy.activity.open"      //系统云端挂机赠送天数活动开关
	AppUinBuyMonthGift      = "buy.month.gift"         //系统云端挂机购买1月赠送天数
	AppUinBuyQuarterGift    = "buy.quarter.gift"       //系统云端挂机购买3月赠送天数
	AppUinBuyHalfGift       = "buy.half.gift"          //系统云端挂机购买6月赠送天数
	AppUinBuyYearGift       = "buy.year.gift"          //系统云端挂机购买12月赠送天数
	AppUinUpdateLevel       = "uin.update.level"       //系统云端挂机最低允许挂机等级
	AppUinUpdateDays        = "uin.update.days"        //系统云端挂机最低允许挂机天数
	AppUinUpdateDate        = "uin.update.date"        //该日期前下单无影响
	AppUinRefundOpen        = "refund.open"            //系统云端挂机退单功能是否开启
	AppUinRefundPoint       = "refund.point"           //系统云端挂机退单功能手续费(百分比)
	AppUinRefundMin         = "refund.min"             //系统云端挂机退单功能手续费(单笔最低)
	AppUinRefundMax         = "refund.max"             //系统云端挂机退单功能手续费(单笔最高)
	AppUinRefundDay         = "refund.day"             //系统云端挂机退单功能最久退单日期
	AppUinRefundFree        = "refund.free"            //系统云端挂机退单功能免费退单日期
	AppUinTaskLevelRealAuth = "task.level.real_auth"   //系统云端挂机n级以下需要钱包实名(-1不限制)
)
