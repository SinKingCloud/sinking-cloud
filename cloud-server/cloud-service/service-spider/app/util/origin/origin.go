package origin

// 结果
const (
	Success     = iota //成功
	Fail               //失败
	TokenError         //token失效
	NotAllow           //token失效
	ProxyExpire        //代理失效
)
