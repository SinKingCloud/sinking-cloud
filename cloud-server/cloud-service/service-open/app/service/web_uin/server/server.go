package server

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

// AppUinServer QQ服务器列表
type AppUinServer struct {
	Key  string `gorm:"column:key" json:"key"`
	Name string `gorm:"column:name" json:"name"`
}
