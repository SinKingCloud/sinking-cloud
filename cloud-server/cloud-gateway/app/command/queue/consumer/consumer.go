package consumer

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

// GetRocketConsumerServer 获取消息队列消费者者对象
func GetRocketConsumerServer(groupName string, retry int) queue.MqConsumer {
	return queue.RegisterRocketConsumerMust([]string{
		util.Consul.Config(group.CloudCommon).Name(config.RocketMqCommon).Viper().GetString(config.RocketmqServer),
	}, groupName, true)
}
