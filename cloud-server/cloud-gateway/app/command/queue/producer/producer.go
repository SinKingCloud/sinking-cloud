package producer

import "github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/queue"

// GetRocketProducerServer 获取消息队列生产者对象
func GetRocketProducerServer(groupName string, retry int) queue.MqProducer {
	return queue.RegisterRocketProducerMust([]string{"42.157.128.40:9876"}, groupName, retry)
}
