package message

import (
	queueMessage "service-user/app/queue/message"
	repositoryMessage "service-user/app/repository/message"
	"service-user/app/service/socket"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
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
	Create(userId int64, title string, content string) error                                                                                                         //创建单条用户消息
	Creates(userIds []int64, title string, content string) error                                                                                                     //批量创建用户消息
	Read(userId int64, ids []int64) error                                                                                                                            //设置用户消息为已读
	ReadAll(userId int64) error                                                                                                                                      //设置用户全部消息为已读
	Delete(userId int64, ids []int64) error                                                                                                                          //删除用户消息
	Clear(userId int64) error                                                                                                                                        //清空用户消息
	SelectByUserId(userId int64, where *repositoryMessage.SelectMessage, orderByField string, orderByType string, page int, pageSize int) ([]*Message, int64, error) //查询用户消息列表，不返回消息内容
	GetInfoByUserId(userId int64, id int64) (*Message, error)                                                                                                        //查询用户消息详情，未读消息自动设置为已读
}

// 消息服务
type service struct {
	repositoryMessage repositoryMessage.Interface
	socketService     socket.Service
	messageQueue      *queue.UnicastClient[*queueMessage.Task]
}

// NewService 实例化service
func NewService(repositoryMessage repositoryMessage.Interface, socketService socket.Service, messageQueue *queue.UnicastClient[*queueMessage.Task]) *service {
	ins := &service{
		repositoryMessage: repositoryMessage,
		socketService:     socketService,
		messageQueue:      messageQueue,
	}
	if messageQueue != nil {
		messageQueue.SetConsumer(ins.consumeCreateTask)
	}
	return ins
}
