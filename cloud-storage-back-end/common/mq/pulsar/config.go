package pulsar

// 生产者配置
type PubConfig struct {
	Enabled bool   // 是否启用 Pulsar
	URL     string // Pulsar 服务地址
	// 主题名称
	Topic string
	// 是否启用批量发送
	BatchingEnabled bool
	// 批量发送最大消息数
	BatchingMaxMessages uint
	// 批量发送最大发布延迟（毫秒）
	BatchingMaxPublishDelay int
	// 压缩类型：None, LZ4, ZLIB, ZSTD
	CompressionType string
	// 发送超时（秒）
	SendTimeout int
	// 是否阻塞队列满时的发送
	BlockIfQueueFull bool
	// 最大等待发送的消息数
	MaxPendingMessages int
}

// 消费者配置
type SubConfig struct {
	Enabled bool   // 是否启用 Pulsar
	URL     string // Pulsar 服务地址
	// 主题名称
	Topic string
	// 多主题订阅
	Topics []string `json:",optional"`
	// 订阅名称
	SubscriptionName string
	// 订阅类型：Exclusive, Shared, Failover, KeyShared
	SubscriptionType string
	// 消费者名称
	Name string `json:",optional"`
	// 初始位置：Latest, Earliest
	SubscriptionInitialPosition string `json:",optional"`
	// 是否自动确认消息
	AutoAck bool `json:",optional"`
	// 未确认消息重新投递延迟（秒）
	NackRedeliveryDelay int `json:",optional"`
	// 接收队列大小
	ReceiverQueueSize int `json:",optional"`
}
