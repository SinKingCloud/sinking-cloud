package server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/constant"
	"strconv"
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
func (s *service) ServerLock(key string) bool {
	return util.Cache.Lock(constant.LockUinServer+key, constant.LockTimeUinServer)
}

// ServerUnLock 服务器解锁
func (s *service) ServerUnLock(key string) bool {
	return util.Cache.UnLock(constant.LockUinServer + key)
}

// ServerIsLock 是否上锁
func (s *service) ServerIsLock(key string) bool {
	return util.Cache.IsLock(constant.LockUinServer + key)
}

// ServerWaitLock 判断服务器是否上锁,上锁的话阻塞任务执行,5s一个判断周期
func (s *service) ServerWaitLock(key string) {
	i := 0
	for true {
		if i >= 120 {
			break
		}
		if util.Cache.IsLock(constant.LockUinServer + key) {
			i++
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}

// ServerInc 自增
func (s *service) ServerInc(key string, second int) int {
	key = constant.LockUinServer + "num_" + key
	if second > 0 {
		defer util.Cache.SetExpire(key, int64(second*1000))
	}
	return util.Cache.Inc(key)
}

// ServerDec 自减
func (s *service) ServerDec(key string, second int) int {
	key = constant.LockUinServer + "num_" + key
	if second > 0 {
		defer util.Cache.SetExpire(key, int64(second*1000))
	}
	return util.Cache.Dec(key)
}

// ServerNum 获取计数器
func (s *service) ServerNum(key string) int {
	str, exists := util.Cache.Get(constant.LockUinServer + "num_" + key)
	if exists && str != "" {
		num, err := strconv.Atoi(str)
		if err == nil {
			return num
		}
		return 0
	}
	return 0
}

// ServerNumWaitZero 等待计数器变为0
func (s *service) ServerNumWaitZero(key string) {
	key = constant.LockUinServer + "num_" + key
	i := 0
	for true {
		if i >= 120 {
			break
		}
		if s.ServerNum(key) > 0 {
			i++
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
