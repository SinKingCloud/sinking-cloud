package queue

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

// GetRocketProducerServer 获取消息队列生产者对象
func GetRocketProducerServer(groupName string, retry int) queue.MqProducer {
	return queue.RegisterRocketProducerMust([]string{
		util.Consul.Config(group.CloudCommon).Name(config.RocketMqCommon).Viper().GetString(config.RocketmqServer),
	}, groupName, retry)
}

// GetRocketConsumerServer 获取消息队列消费者者对象
func GetRocketConsumerServer(groupName string, allListen bool) queue.MqConsumer {
	return queue.RegisterRocketConsumerMust([]string{
		util.Consul.Config(group.CloudCommon).Name(config.RocketMqCommon).Viper().GetString(config.RocketmqServer),
	}, groupName, allListen)
}
