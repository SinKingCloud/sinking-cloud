package uin

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"sync"
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

// AppUin QQ记录
type AppUin struct {
	Id         int          `gorm:"column:id;PRIMARY_KEY" json:"-"`
	UserId     int          `gorm:"column:user_id" json:"user_id"`
	AppId      int          `gorm:"column:app_id" json:"app_id"`
	ServerKey  string       `gorm:"column:server_key" json:"server_key"`
	Uin        string       `gorm:"column:uin" json:"uin"`
	UinStatus  int          `gorm:"column:uin_status" json:"uin_status"`
	SignDay    int          `gorm:"column:sign_day" json:"sign_day"`
	PlayerDay  int          `gorm:"column:player_day" json:"player_day"`
	RunDay     int          `gorm:"column:run_day" json:"run_day"`
	Timing     int          `gorm:"column:timing" json:"timing"`
	Status     int          `gorm:"column:status" json:"status"`
	RunTime    str.DateTime `gorm:"column:run_time" json:"run_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

const (
	Normal = iota //正常
	Stop          //暂停挂机
	Error         //禁止挂机
)

type StatusUin int

const (
	StatusUinNormal StatusUin = iota
	StatusUinError
)
