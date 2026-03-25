package queue

import (
	"encoding/json"
	"errors"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
)

// UnicastClient 单播消息队列客户端
type UnicastClient struct {
	topic        Queue
	name         string
	consumerFunc func(param string) error
	taskChan     chan MqMsg
}

// Try 错误捕获实现
func (l *UnicastClient) Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// SendUnicastTask 发送单播任务
func (l *UnicastClient) SendUnicastTask(param interface{}) {
	marshal, e := json.Marshal(param)
	if e == nil {
		_, _ = l.topic.SendUnicast(l.name, string(marshal))
	}
}

// Stop 停止客户端
func (l *UnicastClient) Stop() error {
	if l.taskChan != nil {
		close(l.taskChan)
	}
	if l.topic != nil {
		l.topic.Close()
	}
	return nil
}

// GetConsumerStatus 获取消费者状态（返回 running，因 worker pool 持续运行）
func (l *UnicastClient) GetConsumerStatus() job.TaskStatus {
	return job.TaskRunning
}

// GetProducerStatus 获取生产者状态（返回 stopped，单播客户端不管理生产者）
func (l *UnicastClient) GetProducerStatus() job.TaskStatus {
	return job.TaskStopped
}

// GetConsumerStats 获取消费者统计
func (l *UnicastClient) GetConsumerStats() job.TaskStats {
	return job.TaskStats{}
}

// GetProducerStats 获取生产者统计（返回空，单播客户端不管理生产者）
func (l *UnicastClient) GetProducerStats() job.TaskStats {
	return job.TaskStats{}
}

// registerUnicast 注册单播消息队列实例
func registerUnicast(queue Queue, topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*UnicastClient, error) {
	if topicName == "" {
		return nil, errors.New("队列名称不能为空")
	}
	if producerThreadNum <= 0 {
		return nil, errors.New("发送者协程数量必须大于0")
	}
	if consumerThreadNum <= 0 {
		return nil, errors.New("消费者协程数量必须大于0")
	}

	ins := &UnicastClient{
		topic:        queue,
		name:         topicName,
		consumerFunc: consumerFunc,
		taskChan:     make(chan MqMsg, consumerThreadNum*100),
	}

	// 预创建 worker pool
	for i := 0; i < consumerThreadNum; i++ {
		go func() {
			for mqMsg := range ins.taskChan {
				ins.Try(func() {
					_ = ins.consumerFunc(mqMsg.Body)
				}, func(e interface{}) {})
			}
		}()
	}

	// 监听单播消息
	err := queue.ListenReceiveUnicast(topicName, func(mqMsg MqMsg) bool {
		ins.taskChan <- mqMsg
		return true
	})
	if err != nil {
		return nil, errors.New("启动" + topicName + "单播消费实例失败: " + err.Error())
	}

	return ins, nil
}

// BroadcastClient 广播消息队列客户端
type BroadcastClient struct {
	topic        Queue
	name         string
	consumerFunc func(param string) error
	taskChan     chan MqMsg
}

// Try 错误捕获实现
func (l *BroadcastClient) Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// SendBroadcastTask 发送广播任务
func (l *BroadcastClient) SendBroadcastTask(param interface{}) {
	marshal, e := json.Marshal(param)
	if e == nil {
		_, _ = l.topic.SendBroadcast(l.name, string(marshal))
	}
}

// Stop 停止客户端
func (l *BroadcastClient) Stop() error {
	if l.taskChan != nil {
		close(l.taskChan)
	}
	if l.topic != nil {
		l.topic.Close()
	}
	return nil
}

// GetConsumerStatus 获取消费者状态（返回 running，因 worker pool 持续运行）
func (l *BroadcastClient) GetConsumerStatus() job.TaskStatus {
	return job.TaskRunning
}

// GetProducerStatus 获取生产者状态（返回 stopped，广播客户端不管理生产者）
func (l *BroadcastClient) GetProducerStatus() job.TaskStatus {
	return job.TaskStopped
}

// GetConsumerStats 获取消费者统计
func (l *BroadcastClient) GetConsumerStats() job.TaskStats {
	return job.TaskStats{}
}

// GetProducerStats 获取生产者统计（返回空，广播客户端不管理生产者）
func (l *BroadcastClient) GetProducerStats() job.TaskStats {
	return job.TaskStats{}
}

// registerBroadcast 注册广播消息队列实例
func registerBroadcast(queue Queue, topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*BroadcastClient, error) {
	if topicName == "" {
		return nil, errors.New("队列名称不能为空")
	}
	if producerThreadNum <= 0 {
		return nil, errors.New("发送者协程数量必须大于0")
	}
	if consumerThreadNum <= 0 {
		return nil, errors.New("消费者协程数量必须大于0")
	}

	ins := &BroadcastClient{
		topic:        queue,
		name:         topicName,
		consumerFunc: consumerFunc,
		taskChan:     make(chan MqMsg, consumerThreadNum*100),
	}

	// 预创建 worker pool
	for i := 0; i < consumerThreadNum; i++ {
		go func() {
			for mqMsg := range ins.taskChan {
				ins.Try(func() {
					_ = ins.consumerFunc(mqMsg.Body)
				}, func(e interface{}) {})
			}
		}()
	}

	// 监听广播消息
	err := queue.ListenReceiveBroadcast(topicName, func(mqMsg MqMsg) bool {
		ins.taskChan <- mqMsg
		return true
	})
	if err != nil {
		return nil, errors.New("启动" + topicName + "广播消费实例失败: " + err.Error())
	}

	return ins, nil
}

// RegisterMemoryUnicast 注册内存单播消息队列实例
func RegisterMemoryUnicast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*UnicastClient, error) {
	topic, err := RegisterMemoryQueueClient()
	if err != nil {
		return nil, errors.New("注册" + topicName + "单播消费实例失败")
	}
	return registerUnicast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterRedisUnicast 注册Redis单播消息队列实例
func RegisterRedisUnicast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*UnicastClient, error) {
	topic, err := RegisterRedisQueueClient(global.App.Redis.GetClient())
	if err != nil {
		return nil, errors.New("注册" + topicName + "单播消费实例失败")
	}
	return registerUnicast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterNatsUnicast 注册NATS单播消息队列实例
func RegisterNatsUnicast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error, natsURL string, options ...interface{}) (*UnicastClient, error) {
	topic, err := RegisterNatsQueueClient(natsURL, options...)
	if err != nil {
		return nil, errors.New("注册" + topicName + "单播消费实例失败: " + err.Error())
	}
	return registerUnicast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterMemoryBroadcast 注册内存广播消息队列实例
func RegisterMemoryBroadcast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*BroadcastClient, error) {
	topic, err := RegisterMemoryQueueClient()
	if err != nil {
		return nil, errors.New("注册" + topicName + "广播消费实例失败")
	}
	return registerBroadcast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterRedisBroadcast 注册Redis广播消息队列实例
func RegisterRedisBroadcast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error) (*BroadcastClient, error) {
	topic, err := RegisterRedisQueueClient(global.App.Redis.GetClient())
	if err != nil {
		return nil, errors.New("注册" + topicName + "广播消费实例失败")
	}
	return registerBroadcast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterNatsBroadcast 注册NATS广播消息队列实例
func RegisterNatsBroadcast(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string) error, natsURL string, options ...interface{}) (*BroadcastClient, error) {
	topic, err := RegisterNatsQueueClient(natsURL, options...)
	if err != nil {
		return nil, errors.New("注册" + topicName + "广播消费实例失败: " + err.Error())
	}
	return registerBroadcast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}
