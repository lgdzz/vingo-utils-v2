package queue

import (
	"fmt"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"github.com/nsqio/go-nsq"
	"sync"
)

type NsqService struct {
	producer     *nsq.Producer
	producerOnce sync.Once
	Addr         string
	Config       *nsq.Config
}

var Nsq NsqService

// 初始化服务（只需要执行1次）
func NsqInit(addr *string, config *nsq.Config) {
	if addr == nil {
		Nsq.Addr = "127.0.0.1:4150"
	}
	if config == nil {
		Nsq.Config = nsq.NewConfig()
		Nsq.Config.MaxInFlight = 1 // 设置每个消费者的最大并发处理消息数为1
		Nsq.Config.MaxAttempts = 0 // 指定重连的最大尝试次数。默认为 0，表示无限次尝试重连。
	}
	Nsq.initProducer()
}

func (s *NsqService) initProducer() {
	var err error
	s.producerOnce.Do(func() {
		s.producer, err = nsq.NewProducer(s.Addr, s.Config)
		if err != nil {
			vingo.LogInfo(fmt.Sprintf("[NSQ]创建生产者失败：%v", err.Error()))
		} else {
			vingo.LogInfo("[NSQ]创建生产者成功.")
		}
	})
}

// ProduceMessageAsync 生产消息（异步）
func (s *NsqService) ProduceMessageAsync(topic string, message []byte) {
	go s.ProduceMessage(topic, message)
}

// ProduceMessage 生产消息
func (s *NsqService) ProduceMessage(topic string, message []byte) {
	err := s.producer.Publish(topic, message)
	if err != nil {
		panic(err.Error())
	}
}

// ConsumeMessagesAsync 消费消息（异步）
func (s *NsqService) ConsumeMessagesAsync(topic string, channel string, handler nsq.Handler) {
	go s.ConsumeMessages(topic, channel, handler)
}

// ConsumeMessages 消费消息
func (s *NsqService) ConsumeMessages(topic string, channel string, handler nsq.Handler) {
	// 创建消费者
	consumer, err := nsq.NewConsumer(topic, channel, s.Config)
	if err != nil {
		vingo.LogInfo(fmt.Sprintf("[NSQ]创建消费者失败：%v", err.Error()))
	} else {
		vingo.LogInfo("[NSQ]创建消费者成功.")
	}

	// 设置消息处理程序
	consumer.AddHandler(handler)

	// 连接到NSQ服务器
	err = consumer.ConnectToNSQD(s.Addr)
	if err != nil {
		vingo.LogInfo(fmt.Sprintf("[NSQ]消费者连接NSQ服务失败：%v", err.Error()))
	}

	// 阻塞等待直到接收到退出信号
	<-consumer.StopChan
}
