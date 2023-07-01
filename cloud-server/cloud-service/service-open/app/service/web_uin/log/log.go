package log

import "sync"

// service 单例对象
type service struct {
}

// obj 单例对象
var (
	obj  *service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *service {
	once.Do(func() {
		obj = &service{}
	})
	return obj
}

const (
	Buy      = iota //在线下单
	Server          //修改地区
	Password        //修改密码
	Login           //登陆账号
	Run             //执行补挂
	Timing          //修改时间
	Delete          //删除挂机
	Status          //修改状态
)
