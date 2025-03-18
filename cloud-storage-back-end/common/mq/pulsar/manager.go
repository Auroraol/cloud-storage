package pulsar

import (
	"context"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/pkg/errors"
)

// PulsarManager 提供 Pulsar 客户端管理功能
type PulsarManager struct {
	client pulsar.Client
	config Config
}

// Config Pulsar 配置
type Config struct {
	// Pulsar 服务地址，例如 "pulsar://localhost:6650"
	URL string
	// 操作超时时间
	OperationTimeout time.Duration
	// 连接超时时间
	ConnectionTimeout time.Duration
	// 是否启用 TLS
	EnableTLS bool
	// TLS 证书路径
	TLSTrustCertsFilePath string
	// 认证类型和参数
	AuthType   string
	AuthParams map[string]string
}

// NewPulsarManager 创建一个新的 Pulsar 管理器
func NewPulsarManager(config Config) (*PulsarManager, error) {
	// 设置默认值
	if config.OperationTimeout == 0 {
		config.OperationTimeout = 30 * time.Second
	}
	if config.ConnectionTimeout == 0 {
		config.ConnectionTimeout = 30 * time.Second
	}

	// 创建客户端选项
	clientOptions := pulsar.ClientOptions{
		URL:               config.URL,
		OperationTimeout:  config.OperationTimeout,
		ConnectionTimeout: config.ConnectionTimeout,
	}

	// 配置 TLS
	if config.EnableTLS {
		clientOptions.TLSTrustCertsFilePath = config.TLSTrustCertsFilePath
		clientOptions.TLSAllowInsecureConnection = false
		clientOptions.TLSValidateHostname = true
	}

	// 配置认证
	if config.AuthType != "" {
		switch config.AuthType {
		case "token":
			token, ok := config.AuthParams["token"]
			if !ok {
				return nil, errors.New("token authentication requires 'token' parameter")
			}
			clientOptions.Authentication = pulsar.NewAuthenticationToken(token)
		case "tls":
			certPath, ok := config.AuthParams["certFile"]
			if !ok {
				return nil, errors.New("TLS authentication requires 'certFile' parameter")
			}
			keyPath, ok := config.AuthParams["keyFile"]
			if !ok {
				return nil, errors.New("TLS authentication requires 'keyFile' parameter")
			}
			clientOptions.Authentication = pulsar.NewAuthenticationTLS(certPath, keyPath)
		default:
			return nil, errors.Errorf("unsupported authentication type: %s", config.AuthType)
		}
	}

	// 创建客户端
	client, err := pulsar.NewClient(clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "创建 Pulsar 客户端失败")
	}

	return &PulsarManager{
		client: client,
		config: config,
	}, nil
}

// Close 关闭 Pulsar 客户端
func (m *PulsarManager) Close() {
	if m.client != nil {
		m.client.Close()
	}
}

// CreateProducer 创建一个生产者
func (m *PulsarManager) CreateProducer(options pulsar.ProducerOptions) (pulsar.Producer, error) {
	if options.Topic == "" {
		return nil, errors.New("topic is required")
	}

	producer, err := m.client.CreateProducer(options)
	if err != nil {
		return nil, errors.Wrap(err, "创建生产者失败")
	}

	return producer, nil
}

// CreateConsumer 创建一个消费者
func (m *PulsarManager) CreateConsumer(options pulsar.ConsumerOptions) (pulsar.Consumer, error) {
	if options.Topic == "" && len(options.Topics) == 0 {
		return nil, errors.New("topic or topics is required")
	}

	if options.SubscriptionName == "" {
		return nil, errors.New("subscription name is required")
	}

	consumer, err := m.client.Subscribe(options)
	if err != nil {
		return nil, errors.Wrap(err, "创建消费者失败")
	}

	return consumer, nil
}

// GetClient 获取原始 Pulsar 客户端
func (m *PulsarManager) GetClient() pulsar.Client {
	return m.client
}

// SendMessage 发送消息到指定主题
func (m *PulsarManager) SendMessage(ctx context.Context, topic string, payload []byte, properties map[string]string) error {
	producer, err := m.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return err
	}
	defer producer.Close()

	msg := pulsar.ProducerMessage{
		Payload:    payload,
		Properties: properties,
	}

	_, err = producer.Send(ctx, &msg)
	return err
}

// SendMessageAsync 异步发送消息到指定主题
func (m *PulsarManager) SendMessageAsync(ctx context.Context, topic string, payload []byte, properties map[string]string, callback func(pulsar.MessageID, *pulsar.ProducerMessage, error)) {
	producer, err := m.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		if callback != nil {
			callback(nil, nil, err)
		}
		return
	}
	defer producer.Close()

	msg := pulsar.ProducerMessage{
		Payload:    payload,
		Properties: properties,
	}

	producer.SendAsync(ctx, &msg, callback)
}
