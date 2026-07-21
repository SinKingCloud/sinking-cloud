package mind

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// Data 导图数据接口返回数据，value 为空时不返回该字段
type Data struct {
	Id         int64        `json:"id"`
	TypeId     int64        `json:"type_id"`
	Title      string       `json:"title"`
	Describe   string       `json:"describe"`
	Value      string       `json:"value,omitempty"`
	Shared     bool         `json:"shared"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}

// Type 导图分类接口返回数据
type Type struct {
	Id         int64        `json:"id"`
	Name       string       `json:"name"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}

// Version 导图版本接口返回数据
type Version struct {
	Id         int64        `json:"id"`
	DataId     int64        `json:"data_id"`
	Type       int          `json:"type"`
	Value      string       `json:"value,omitempty"`
	Remark     string       `json:"remark"`
	Current    bool         `json:"current"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}

// SharedData 公开分享的导图数据
type SharedData struct {
	Title      string       `json:"title"`
	Describe   string       `json:"describe"`
	Value      string       `json:"value"`
	UpdateTime str.DateTime `json:"update_time"`
}
