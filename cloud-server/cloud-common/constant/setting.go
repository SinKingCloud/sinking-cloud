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

	UiGroup     = "ui"                   //界面配置
	UiLayout    = UiGroup + ".layout"    //界面布局
	UiWaterMark = UiGroup + ".watermark" //界面水印
	UiTheme     = UiGroup + ".theme"     //界面主题
	UiCompact   = UiGroup + ".compact"   //紧凑模式
	UiColor     = UiGroup + ".color"     //主题色
	UiRadius    = UiGroup + ".radius"    //主题圆角

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
