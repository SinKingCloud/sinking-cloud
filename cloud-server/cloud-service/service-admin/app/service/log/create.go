package log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud/service/service-admin/app/command/queues/producer"
)

// Create 插入数据
func (s *service) Create(data *model.SystemLog) (err error) {
	err = util.Database.Db.Create(&data).Error
	return
}

// CreateByQueue 插入数据(异步消息队列)
func (s *service) CreateByQueue(userId int, requestId string, requestIp string, types int, title string, content string) {
	producer.NewLogProducer().Build(&model.SystemLog{
		UserId:     userId,
		RequestId:  requestId,
		RequestIp:  requestIp,
		Type:       types,
		Title:      title,
		Content:    content,
		CreateTime: str.DateTime{},
		UpdateTime: str.DateTime{},
	}).Send()
}
