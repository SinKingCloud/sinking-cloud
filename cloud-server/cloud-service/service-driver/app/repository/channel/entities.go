package channel

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// Channel 存储通道
type Channel struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Type       int          `gorm:"column:type" json:"type"`
	Name       string       `gorm:"column:name" json:"name"`
	Status     int          `gorm:"column:status" json:"status"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

type SelectChannel struct {
	Name            *string
	Type            *int
	Status          *int
	CreateTimeStart *string
	CreateTimeEnd   *string
	UpdateTimeStart *string
	UpdateTimeEnd   *string
}

type UpdateChannel struct {
	Type   *int
	Name   *string
	Config *string
	Status *int
}
