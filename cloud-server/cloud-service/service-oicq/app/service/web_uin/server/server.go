package server

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-oicq/app/constant"
	"strconv"
	"sync"
	"time"
)

// Service 单例对象
type Service struct {
}

// obj 单例对象
var (
	obj  *Service
	once sync.Once
)

// GetIns 获取单例
func GetIns() *Service {
	once.Do(func() {
		obj = &Service{}
	})
	return obj
}

// AppUinServer QQ服务器列表
type AppUinServer struct {
	Key  string `gorm:"column:key" json:"key"`
	Name string `gorm:"column:name" json:"name"`
}

// ServerLock 服务器上锁。阻塞挂机执行
func (s *Service) ServerLock(key string) bool {
	return util.Cache.Lock(constant.LockUinServer+key, constant.LockTimeUinServer)
}

// ServerUnLock 服务器解锁
func (s *Service) ServerUnLock(key string) bool {
	return util.Cache.UnLock(constant.LockUinServer + key)
}

// ServerIsLock 是否上锁
func (s *Service) ServerIsLock(key string) bool {
	return util.Cache.IsLock(constant.LockUinServer + key)
}

// ServerWaitLock 判断服务器是否上锁,上锁的话阻塞任务执行,5s一个判断周期 是否可以继续执行
func (s *Service) ServerWaitLock(key string) bool {
	res := false
	i := 0
	for true {
		if i >= 720 {
			break
		}
		if util.Cache.IsLock(constant.LockUinServer + key) {
			i++
			time.Sleep(5 * time.Second)
		} else {
			res = true
			break
		}
	}
	return res
}

// ServerInc 自增
func (s *Service) ServerInc(key string, second int) int {
	key = constant.LockUinServer + "num_" + key
	if second > 0 {
		defer util.Cache.SetExpire(key, int64(second*1000))
	}
	return util.Cache.Inc(key)
}

// ServerDec 自减
func (s *Service) ServerDec(key string, second int) int {
	key = constant.LockUinServer + "num_" + key
	if second > 0 {
		defer util.Cache.SetExpire(key, int64(second*1000))
	}
	return util.Cache.Dec(key)
}

// ServerReset 地区重置
func (s *Service) ServerReset(key string) bool {
	key = constant.LockUinServer + "num_" + key
	return util.Cache.Delete(key)
}

// ServerJobNumInc 地区并发任务自增
func (s *Service) ServerJobNumInc(key string) int {
	key = constant.CacheNameWithUinServerJob + key
	defer util.Cache.SetExpire(key, int64(constant.CacheTimeWithCacheNameUinServerJob))
	return util.Cache.Inc(key)
}

// ServerJobNumDec 地区并发任务自减
func (s *Service) ServerJobNumDec(key string) int {
	key = constant.CacheNameWithUinServerJob + key
	defer util.Cache.SetExpire(key, int64(constant.CacheTimeWithCacheNameUinServerJob))
	return util.Cache.Dec(key)
}

// ServerJobReset 地区并发任务重置
func (s *Service) ServerJobReset(key string) bool {
	key = constant.CacheNameWithUinServerJob + key
	return util.Cache.Delete(key)
}

// ServerNum 获取计数器
func (s *Service) ServerNum(key string) int {
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

// ServerJobNum 获取地区任务数量
func (s *Service) ServerJobNum(key string) int {
	str, exists := util.Cache.Get(constant.CacheNameWithUinServerJob + key)
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
func (s *Service) ServerNumWaitZero(key string) {
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
