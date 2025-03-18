package pulsar

import (
	"context"
	"encoding/json"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/pkg/errors"
)

// Publisher Pulsar 消息发布者
type Publisher struct {
	manager  *PulsarManager
	producer pulsar.Producer
	topic    string
	config   PubConfig
}

// NewPublisher 创建一个新的消息发布者
func NewPublisher(manager *PulsarManager, config PubConfig) (*Publisher, error) {
	if config.Topic == "" {
		return nil, errors.New("topic is required")
	}

	// 创建生产者选项
	options := pulsar.ProducerOptions{
		Topic:              config.Topic,
		MaxPendingMessages: config.MaxPendingMessages,
	}

	// 设置批量发送
	if config.BatchingEnabled {
		options.BatchingMaxMessages = config.BatchingMaxMessages
		options.BatchingMaxPublishDelay = time.Duration(config.BatchingMaxPublishDelay) * time.Millisecond
	} else {
		options.BatchingMaxMessages = 1
	}

	// 设置压缩类型
	switch config.CompressionType {
	case "LZ4":
		options.CompressionType = pulsar.LZ4
	case "ZLIB":
		options.CompressionType = pulsar.ZLib
	case "ZSTD":
		options.CompressionType = pulsar.ZSTD
	default:
		options.CompressionType = pulsar.NoCompression
	}

	// 设置发送超时
	if config.SendTimeout > 0 {
		options.SendTimeout = time.Duration(config.SendTimeout) * time.Second
	}

	// 创建生产者
	producer, err := manager.CreateProducer(options)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		manager:  manager,
		producer: producer,
		topic:    config.Topic,
		config:   config,
	}, nil
}

// Close 关闭生产者
func (p *Publisher) Close() {
	if p.producer != nil {
		p.producer.Close()
	}
}

// Send 发送消息
func (p *Publisher) Send(ctx context.Context, payload []byte, properties map[string]string) (pulsar.MessageID, error) {
	msg := pulsar.ProducerMessage{
		Payload:    payload,
		Properties: properties,
	}

	return p.producer.Send(ctx, &msg)
}

// SendAsync 异步发送消息
func (p *Publisher) SendAsync(ctx context.Context, payload []byte, properties map[string]string, callback func(pulsar.MessageID, *pulsar.ProducerMessage, error)) {
	msg := pulsar.ProducerMessage{
		Payload:    payload,
		Properties: properties,
	}

	p.producer.SendAsync(ctx, &msg, callback)
}

// SendObject 发送对象（将对象序列化为 JSON）
func (p *Publisher) SendObject(ctx context.Context, obj interface{}, properties map[string]string) (pulsar.MessageID, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, errors.Wrap(err, "序列化对象失败")
	}

	return p.Send(ctx, data, properties)
}

// SendObjectAsync 异步发送对象（将对象序列化为 JSON）
func (p *Publisher) SendObjectAsync(ctx context.Context, obj interface{}, properties map[string]string, callback func(pulsar.MessageID, *pulsar.ProducerMessage, error)) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "序列化对象失败")
	}

	p.SendAsync(ctx, data, properties, callback)
	return nil
}

// GetTopic 获取主题
func (p *Publisher) GetTopic() string {
	return p.topic
}

// GetProducer 获取原始生产者
func (p *Publisher) GetProducer() pulsar.Producer {
	return p.producer
}
