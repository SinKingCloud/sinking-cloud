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
	Id         int64        `json:"id"`
	Type       int          `json:"type"`
	Name       string       `json:"name"`
	Spec       string       `json:"spec"`
	Concurrent int          `json:"concurrent"`
	Status     int          `json:"status"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}
