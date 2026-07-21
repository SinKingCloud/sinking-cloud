package mind_version

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// MindVersion 导图版本列表展示数据，不包含 value 字段
type MindVersion struct {
	Id         int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	DataId     int64        `gorm:"column:data_id" json:"data_id"`
	Type       int          `gorm:"column:type" json:"type"`
	Remark     string       `gorm:"column:remark" json:"remark"`
	CreateTime str.DateTime `gorm:"column:create_time" json:"create_time"`
	UpdateTime str.DateTime `gorm:"column:update_time" json:"update_time"`
}

// SelectMindVersion 导图版本查询
type SelectMindVersion struct {
	DataId int64
	Type   string
}

// UpdateMindVersion 导图版本更新
type UpdateMindVersion struct {
	Value  string
	Remark interface{}
}
