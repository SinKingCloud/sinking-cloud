package order_log

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"

// OrderLog 订单日志表
type OrderLog struct {
	Id           int64        `gorm:"column:id;PRIMARY_KEY" json:"id"`
	OrderId      int64        `gorm:"column:order_id" json:"order_id"`
	OperatorId   int64        `gorm:"column:operator_id" json:"operator_id"`
	OperatorType int64        `gorm:"column:operator_type" json:"operator_type"`
	Type         int          `gorm:"column:type" json:"type"`
	Content      int          `gorm:"column:content" json:"content"`
	CreateTime   str.DateTime `gorm:"column:create_time" json:"create_time"`
}

type SelectOrderLog struct {
	OrderId         string
	OperatorId      string
	OperatorType    string
	Type            string
	CreateTimeStart string
	CreateTimeEnd   string
}
