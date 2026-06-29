package message

import (
	repositoryMessage "service-user/app/repository/message"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
)

// Message 消息接口返回数据，隐藏 user_id；content 为空时不返回该字段。
type Message struct {
	Id         int64        `json:"id"`
	Title      string       `json:"title"`
	Content    string       `json:"content,omitempty"`
	Status     int          `json:"status"`
	ReadTime   str.DateTime `json:"read_time"`
	CreateTime str.DateTime `json:"create_time"`
	UpdateTime str.DateTime `json:"update_time"`
}

// Service 消息服务接口
type Service interface {
	Create(userId int64, title string, content string) error
	Creates(userIds []int64, title string, content string) error
	Read(userId int64, ids []int64) error
	Delete(userId int64, ids []int64) error
	Clear(userId int64) error
	SelectByUserId(userId int64, where *repositoryMessage.SelectMessage, orderByField string, orderByType string, page int, pageSize int) ([]*Message, int64, error)
	GetInfoByUserId(userId int64, id int64) (*Message, error)
}

// 消息服务
type service struct {
	repositoryMessage repositoryMessage.Interface
}

// NewService 实例化service
func NewService(repositoryMessage repositoryMessage.Interface) *service {
	return &service{
		repositoryMessage: repositoryMessage,
	}
}
