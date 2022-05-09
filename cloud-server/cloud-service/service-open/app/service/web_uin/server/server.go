package server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"time"
)

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

// ServerLock 服务器上锁。阻塞挂机执行
func (s *service) ServerLock(addr string) bool {
	return util.Cache.Lock(constant.LockUinServer+addr, constant.LockTimeUinServer)
}

// ServerUnLock 服务器服务器解锁
func (s *service) ServerUnLock(addr string) bool {
	return util.Cache.UnLock(constant.LockUinServer + addr)
}

// ServerWaitLock 判断服务器是否上锁,上锁的话阻塞执行
func (s *service) ServerWaitLock(addr string) {
	for true {
		if util.Cache.IsLock(constant.LockUinServer + addr) {
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}

const (
	Normal = iota //正常
	Error         //禁用
)
