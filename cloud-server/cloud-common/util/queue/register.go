package queue

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/global"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/job"
)

type Client struct {
	topic             Queue
	name              string
	consumer          chan interface{}
	consumerThreadNum int
	consumerFunc      func(param string)
	producer          chan interface{}
	producerThreadNum int

	// 任务管理
	consumerTask *job.Task
	producerTask *job.Task
}

// Try 错误捕获实现
func (l *Client) Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func (l *Client) run() error {
	if l.consumer == nil {
		l.consumer = make(chan interface{}, 10000)
	}
	if l.producer == nil {
		l.producer = make(chan interface{}, 10000)
	}

	// 启动消费者任务
	l.consumerTask = job.NewTask().
		SetThread(l.consumerThreadNum).
		SetProducer(func(ctx context.Context, channel chan interface{}) error {
			l.consumer = channel
			<-ctx.Done() // 等待上下文取消
			return ctx.Err()
		}).
		SetConsumer(func(ctx context.Context, param interface{}) error {
			l.Try(func() {
				l.consumerFunc(param.(string))
			}, nil)
			return nil
		})

	if err := l.consumerTask.Start(); err != nil {
		return err
	}

	// 启动生产者任务
	l.producerTask = job.NewTask().
		SetThread(l.producerThreadNum).
		SetProducer(func(ctx context.Context, channel chan interface{}) error {
			l.producer = channel
			<-ctx.Done() // 等待上下文取消
			return ctx.Err()
		}).
		SetConsumer(func(ctx context.Context, param interface{}) error {
			l.Try(func() {
				marshal, e := json.Marshal(param)
				if e == nil {
					_, e = l.topic.SendMsg(l.name, string(marshal))
				}
			}, nil)
			return nil
		})

	if err := l.producerTask.Start(); err != nil {
		return err
	}

	err := l.topic.ListenReceiveMsgDo(l.name, func(mqMsg MqMsg) {
		l.consumer <- mqMsg.Body
	})
	if err != nil {
		return err
	}
	return nil
}

// SendTask 发送任务
func (l *Client) SendTask(param interface{}) {
	l.producer <- param
}

// Stop 停止客户端
func (l *Client) Stop() error {
	var err error
	if l.consumerTask != nil {
		if stopErr := l.consumerTask.Stop(); stopErr != nil {
			err = stopErr
		}
	}
	if l.producerTask != nil {
		if stopErr := l.producerTask.Stop(); stopErr != nil {
			err = stopErr
		}
	}

	// 关闭底层队列连接
	if l.topic != nil {
		l.topic.Close()
	}

	return err
}

// GetConsumerStatus 获取消费者任务状态
func (l *Client) GetConsumerStatus() job.TaskStatus {
	if l.consumerTask == nil {
		return job.TaskStopped
	}
	return l.consumerTask.GetStatus()
}

// GetProducerStatus 获取生产者任务状态
func (l *Client) GetProducerStatus() job.TaskStatus {
	if l.producerTask == nil {
		return job.TaskStopped
	}
	return l.producerTask.GetStatus()
}

// GetConsumerStats 获取消费者任务统计信息
func (l *Client) GetConsumerStats() job.TaskStats {
	if l.consumerTask == nil {
		return job.TaskStats{}
	}
	return l.consumerTask.GetStats()
}

// GetProducerStats 获取生产者任务统计信息
func (l *Client) GetProducerStats() job.TaskStats {
	if l.producerTask == nil {
		return job.TaskStats{}
	}
	return l.producerTask.GetStats()
}

// register 注册消息队列实例
// topicName 队列名称
// producerThreadNum 发送者协程数量
// consumerThreadNum 消费者协程数量
func register(queue Queue, topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string)) (*Client, error) {
	if topicName == "" {
		return nil, errors.New("队列名称不能为空")
	}
	if producerThreadNum <= 0 {
		return nil, errors.New("发送者协程数量必须大于0")
	}
	if consumerThreadNum <= 0 {
		return nil, errors.New("消费者协程数量必须大于0")
	}
	ins := &Client{
		topic:             queue,
		name:              topicName,
		producerThreadNum: producerThreadNum,
		consumerThreadNum: consumerThreadNum,
		consumerFunc:      consumerFunc,
	}
	err := ins.run()
	if err != nil {
		return nil, errors.New("启动" + topicName + "消费实例失败")
	}
	return ins, nil
}

// RegisterMemory 注册redis消息队列实例
// topicName 队列名称
// producerThreadNum 发送者协程数量
// consumerThreadNum 消费者协程数量
func RegisterMemory(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string)) (*Client, error) {
	topic, err := RegisterMemoryQueueClient(100) //注册消费实例
	if err != nil {
		return nil, errors.New("注册" + topicName + "消费实例失败")
	}
	return register(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterRedis 注册redis消息队列实例
// topicName 队列名称
// producerThreadNum 发送者协程数量
// consumerThreadNum 消费者协程数量
func RegisterRedis(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string)) (*Client, error) {
	topic, err := RegisterRedisQueueClient(global.Redis.GetClient()) //注册消费实例
	if err != nil {
		return nil, errors.New("注册" + topicName + "消费实例失败")
	}
	return register(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterNats 注册NATS消息队列实例
// topicName 队列名称
// producerThreadNum 发送者协程数量
// consumerThreadNum 消费者协程数量
// natsURL NATS服务器地址 (如: nats://localhost:4222)
// options NATS连接选项 (可选)
func RegisterNats(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string), natsURL string, options ...interface{}) (*Client, error) {
	topic, err := RegisterNatsQueueClient(natsURL, options...) //注册消费实例
	if err != nil {
		return nil, errors.New("注册" + topicName + "消费实例失败: " + err.Error())
	}
	return register(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}
