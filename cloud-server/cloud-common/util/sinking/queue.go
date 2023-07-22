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
	vip := util.Consul.Config(group.CloudCommon).Name(config.RocketMqCommon).Viper()
	return queue.RegisterRocketProducerMust([]string{
		vip.GetString(config.RocketmqServer),
	}, groupName, retry, vip.GetString(config.RocketmqAccessKey), vip.GetString(config.RocketmqSecretKey))
}

// GetRocketConsumerServer 获取消息队列消费者者对象
func GetRocketConsumerServer(groupName string, allListen bool) queue.MqConsumer {
	vip := util.Consul.Config(group.CloudCommon).Name(config.RocketMqCommon).Viper()
	return queue.RegisterRocketConsumerMust([]string{
		vip.GetString(config.RocketmqServer),
	}, groupName, allListen, vip.GetString(config.RocketmqAccessKey), vip.GetString(config.RocketmqSecretKey))
}

// GetServerHash 获取节点hash
func GetServerHash(name string) string {
	return name + str.NewStringTool(util.Conf.GetString(config.ConfApplicationAddr)+":"+util.Conf.GetString(config.ConfApplicationPort)).Md5()
}
