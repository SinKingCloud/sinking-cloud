package system_log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-admin/app/command/queues/producer"
)

// Create 插入数据
func (s *Service) Create(data *model.SystemLog) (err error) {
	err = util.Database.Db.Create(&data).Error
	return
}

// CreateByQueue 插入数据(异步消息队列)
func (s *Service) CreateByQueue(c *sinking_server.Context, types int, title string, content string) {
	producer.NewLogProducer().Build(&model.SystemLog{
		UserId:     c.GetSystemUserInfo().Id,
		RequestId:  c.GetRequestId(),
		RequestIp:  c.GetRequestIp(),
		Type:       types,
		Title:      title,
		Content:    content,
		CreateTime: str.DateTime{},
		UpdateTime: str.DateTime{},
	}).Send()
}
