package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"go.uber.org/zap"
)

// PulsarConfig Pulsar配置
type PulsarConfig struct {
	URL               string        // Pulsar服务地址
	OperationTimeout  time.Duration // 操作超时时间
	ConnectionTimeout time.Duration // 连接超时时间
}

// PulsarManager Pulsar客户端管理器
type PulsarManager struct {
	client pulsar.Client
	config PulsarConfig
	mu     sync.RWMutex
	// 缓存已创建的生产者
	producers map[string]pulsar.Producer
}

// NewPulsarManager 创建Pulsar客户端管理器
func NewPulsarManager(config PulsarConfig) (*PulsarManager, error) {
	// 设置默认值
	if config.OperationTimeout == 0 {
		config.OperationTimeout = 30 * time.Second
	}
	if config.ConnectionTimeout == 0 {
		config.ConnectionTimeout = 30 * time.Second
	}

	// 创建Pulsar客户端
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               config.URL,
		OperationTimeout:  config.OperationTimeout,
		ConnectionTimeout: config.ConnectionTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("创建Pulsar客户端失败: %v", err)
	}

	return &PulsarManager{
		client:    client,
		config:    config,
		producers: make(map[string]pulsar.Producer),
	}, nil
}

// Close 关闭Pulsar客户端
func (m *PulsarManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 关闭所有生产者
	for _, producer := range m.producers {
		producer.Close()
	}

	// 关闭客户端
	m.client.Close()
}

// GetProducer 获取或创建生产者
func (m *PulsarManager) GetProducer(topic string) (pulsar.Producer, error) {
	m.mu.RLock()
	producer, ok := m.producers[topic]
	m.mu.RUnlock()

	if ok {
		return producer, nil
	}

	// 如果不存在，创建新的生产者
	m.mu.Lock()
	defer m.mu.Unlock()

	// 再次检查，避免并发创建
	producer, ok = m.producers[topic]
	if ok {
		return producer, nil
	}

	// 创建生产者
	producer, err := m.client.CreateProducer(pulsar.ProducerOptions{
		Topic:           topic,
		SendTimeout:     m.config.OperationTimeout,
		DisableBatching: false,
	})
	if err != nil {
		return nil, fmt.Errorf("创建生产者失败: %v", err)
	}

	// 缓存生产者
	m.producers[topic] = producer
	return producer, nil
}

// SendMessage 发送消息
func (m *PulsarManager) SendMessage(topic string, payload []byte, properties map[string]string) (pulsar.MessageID, error) {
	producer, err := m.GetProducer(topic)
	if err != nil {
		return nil, err
	}

	// 创建消息
	msg := &pulsar.ProducerMessage{
		Payload:    payload,
		Properties: properties,
		EventTime:  time.Now(),
	}

	// 发送消息
	return producer.Send(context.Background(), msg)
}

// SendMessageAsync 异步发送消息
func (m *PulsarManager) SendMessageAsync(topic string, payload []byte, properties map[string]string, callback func(pulsar.MessageID, *pulsar.ProducerMessage, error)) {
	producer, err := m.GetProducer(topic)
	if err != nil {
		if callback != nil {
			callback(nil, nil, err)
		}
		return
	}

	// 创建消息
	msg := &pulsar.ProducerMessage{
		Payload:    payload,
		Properties: properties,
		EventTime:  time.Now(),
	}

	// 异步发送消息
	producer.SendAsync(context.Background(), msg, callback)
}

// CreateConsumer 创建消费者
func (m *PulsarManager) CreateConsumer(options pulsar.ConsumerOptions) (pulsar.Consumer, error) {
	// 创建消费者
	return m.client.Subscribe(options)
}

// ConsumeMessages 消费消息
func (m *PulsarManager) ConsumeMessages(consumer pulsar.Consumer, handler func(pulsar.Message) error) {
	for {
		// 接收消息
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			zap.S().Errorf("接收消息失败: %v", err)
			time.Sleep(1 * time.Second) // 避免CPU空转
			continue
		}

		// 处理消息
		err = handler(msg)
		if err != nil {
			zap.S().Errorf("处理消息失败: %v", err)
			consumer.Nack(msg)
		} else {
			consumer.Ack(msg)
		}
	}
}
