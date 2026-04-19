package task

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// SelectTask 任务查询
type SelectTask struct {
	Type            string
	Name            string
	Status          string
	Concurrent      string
	CreateTimeStart string
	CreateTimeEnd   string
	UpdateTimeStart string
	UpdateTimeEnd   string
}

// UpdateTask 任务更新
type UpdateTask struct {
	Type       interface{}
	Name       interface{}
	Spec       interface{}
	Script     interface{}
	Concurrent interface{}
	Status     interface{}
}

// Task 计划任务
type Task struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Type       int          `gorm:"column:type" json:"type"`
	Name       string       `gorm:"column:name" json:"name"`
	Spec       string       `gorm:"column:spec" json:"spec"`
	Concurrent int          `gorm:"column:concurrent" json:"concurrent"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}
