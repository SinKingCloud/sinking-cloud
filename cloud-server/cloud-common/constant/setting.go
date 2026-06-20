package constant

const (
	WebGroup    = "web"                  //网站配置
	WebUrl      = WebGroup + ".url"      //网站地址
	WebApi      = WebGroup + ".api"      //网站API地址
	WebName     = WebGroup + ".name"     //网站名称
	WebTitle    = WebGroup + ".title"    //网站标题
	WebKeyWords = WebGroup + ".keywords" //网站关键词
	WebDescribe = WebGroup + ".describe" //网站描述
	WebContact  = WebGroup + ".contact"  //网站联系人
	WebLogo     = WebGroup + ".logo"     //网站logo

	UiGroup          = "ui"                         //界面配置
	UiUserLayout     = UiGroup + ".user.layout"     //前台界面布局
	UiUserWaterMark  = UiGroup + ".user.watermark"  //前台界面水印
	UiUserTheme      = UiGroup + ".user.theme"      //前台界面主题
	UiUserCompact    = UiGroup + ".user.compact"    //前台紧凑模式
	UiUserColor      = UiGroup + ".user.color"      //前台主题色
	UiUserRadius     = UiGroup + ".user.radius"     //前台主题圆角
	UiAdminLayout    = UiGroup + ".admin.layout"    //后台界面布局
	UiAdminWaterMark = UiGroup + ".admin.watermark" //后台界面水印
	UiAdminTheme     = UiGroup + ".admin.theme"     //后台界面主题
	UiAdminCompact   = UiGroup + ".admin.compact"   //后台紧凑模式
	UiAdminColor     = UiGroup + ".admin.color"     //后台主题色
	UiAdminRadius    = UiGroup + ".admin.radius"    //后台主题圆角

	AliGroup       = "ali"                      //阿里云配置
	AliKey         = AliGroup + ".key"          //阿里云key
	AliSecret      = AliGroup + ".secret"       //阿里云secret
	AliSmsSign     = AliGroup + ".sms.sign"     //阿里云sms sign
	AliSmsTemplate = AliGroup + ".sms.template" //阿里云sms template
	AliSmsVar      = AliGroup + ".sms.var"      //阿里云sms变量名称
	AliMailAccount = AliGroup + ".mail.account" //阿里云mail 发信地址
	AliFaceScene   = AliGroup + ".face.scene"   //阿里云实名认证场景ID
	AliFaceDomain  = AliGroup + ".face.domain"  //阿里云实名认证场景资料储存oss域名
)
