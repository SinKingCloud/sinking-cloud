package producer

import (
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/config"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/constant/service/group"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"
)

// GetRocketProducerServer 获取消息队列生产者对象
func GetRocketProducerServer(groupName string, retry int) queue.MqProducer {
	fmt.Println(util.Consul.Config(group.CloudCommon).Name(config.RocketMqCommon).Viper().GetString(config.RocketmqServer))
	return queue.RegisterRocketProducerMust([]string{
		util.Consul.Config(group.CloudCommon).Name(config.RocketMqCommon).Viper().GetString(config.RocketmqServer),
	}, groupName, retry)
}
