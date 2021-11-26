package sinking

import (
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
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

// GetServerHash 获取节点hash
func GetServerHash(name string) string {
	return name + str.NetStringTool(util.Conf.GetString(config.ConfApplicationAddr)+":"+util.Conf.GetString(config.ConfApplicationPort)).Md5()
}
