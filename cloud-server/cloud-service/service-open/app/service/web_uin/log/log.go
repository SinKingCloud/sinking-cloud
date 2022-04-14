package log

// service 单例对象
type service struct {
}

// obj 单例对象
var obj *service

// GetIns 获取单例
func GetIns() *service {
	if obj == nil {
		obj = &service{}
	}
	return obj
}

const (
	Buy      = iota //在线下单
	Server          //修改地区
	Password        //修改密码
	Login           //登陆账号
	Run             //执行补挂
)
