package log

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/sinking_server"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-user/app/command/queues/producer"
)

// Create 插入数据
func (s *service) Create(data *model.SystemLog) (err error) {
	err = util.Database.Db.Create(&data).Error
	return
}

// CreateByQueue 插入数据(异步消息队列)
func (s *service) CreateByQueue(c *sinking_server.Context, types int, title string, content string) {
	producer.NewWebUserLogProducer().Build(&model.WebUserLog{
		UserId:    c.GetWebUserInfo().Id,
		RequestId: c.GetRequestId(),
		RequestIp: c.GetRequestIp(),
		Type:      types,
		Title:     title,
		Content:   content,
	}).Send()
}
