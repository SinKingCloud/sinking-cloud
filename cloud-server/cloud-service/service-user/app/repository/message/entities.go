package message

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// SelectMessage 消息查询
type SelectMessage struct {
	UserId          *int64
	Title           *string
	Status          *int
	ReadTimeStart   *string
	ReadTimeEnd     *string
	CreateTimeStart *string
	CreateTimeEnd   *string
	UpdateTimeStart *string
	UpdateTimeEnd   *string
}

// UpdateMessage 消息更新
type UpdateMessage struct {
	Status   *int
	ReadTime *str.DateTime
}

// Message 消息列表展示数据，不包含 content 字段。
type Message struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	Title      string       `gorm:"column:title" json:"title"`
	Status     int          `gorm:"column:status" json:"status"`
	ReadTime   str.DateTime `gorm:"column:read_time" json:"read_time"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}
