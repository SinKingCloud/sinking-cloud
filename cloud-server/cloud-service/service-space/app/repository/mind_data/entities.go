package mind_data

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// MindData 导图数据列表展示数据，不包含 value 字段
type MindData struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	UserId     int64        `gorm:"column:user_id" json:"user_id"`
	TypeId     int64        `gorm:"column:type_id" json:"type_id"`
	Title      string       `gorm:"column:title" json:"title"`
	Describe   string       `gorm:"column:describe" json:"describe"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// SelectMindData 导图数据查询
type SelectMindData struct {
	UserId  *int64
	TypeId  *int64
	Keyword *string
}

// UpdateMindData 导图数据更新
type UpdateMindData struct {
	TypeId   *int64
	Title    *string
	Describe *string
	Value    *string
}
