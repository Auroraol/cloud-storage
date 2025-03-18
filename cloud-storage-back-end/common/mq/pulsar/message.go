package pulsar

import (
	"time"
)

// BaseMessage 基础消息结构
type BaseMessage struct {
	// 消息ID
	MessageID string `json:"message_id"`
	// 消息类型
	MessageType string `json:"message_type"`
	// 消息发送时间
	SendTime time.Time `json:"send_time"`
	// 消息来源
	Source string `json:"source"`
	// 消息版本
	Version string `json:"version"`
}

// FileUploadedMessage 文件上传完成消息
type FileUploadedMessage struct {
	BaseMessage
	// 文件ID
	FileID string `json:"file_id"`
	// 文件名
	FileName string `json:"file_name"`
	// 文件大小（字节）
	FileSize int64 `json:"file_size"`
	// 文件内容类型
	ContentType string `json:"content_type"`
	// 用户ID
	UserID string `json:"user_id"`
	// 上传时间
	UploadTime time.Time `json:"upload_time"`
	// 存储路径
	StoragePath string `json:"storage_path"`
	// 文件哈希
	FileHash string `json:"file_hash,omitempty"`
	// 文件元数据
	Metadata map[string]string `json:"metadata,omitempty"`
}

// FileProcessedMessage 文件处理完成消息
type FileProcessedMessage struct {
	BaseMessage
	// 文件ID
	FileID string `json:"file_id"`
	// 处理类型（如：压缩、转码、缩略图等）
	ProcessType string `json:"process_type"`
	// 处理结果
	Result string `json:"result"`
	// 处理时间
	ProcessTime time.Time `json:"process_time"`
	// 处理后的文件路径
	OutputPath string `json:"output_path,omitempty"`
	// 处理耗时（毫秒）
	Duration int64 `json:"duration"`
	// 处理状态（成功/失败）
	Status string `json:"status"`
	// 错误信息（如果处理失败）
	ErrorMessage string `json:"error_message,omitempty"`
}

// UserActivityMessage 用户活动消息
type UserActivityMessage struct {
	BaseMessage
	// 用户ID
	UserID string `json:"user_id"`
	// 活动类型
	ActivityType string `json:"activity_type"`
	// 活动时间
	ActivityTime time.Time `json:"activity_time"`
	// 活动详情
	Details map[string]interface{} `json:"details,omitempty"`
	// IP地址
	IPAddress string `json:"ip_address,omitempty"`
	// 用户代理
	UserAgent string `json:"user_agent,omitempty"`
}

// SystemNotificationMessage 系统通知消息
type SystemNotificationMessage struct {
	BaseMessage
	// 通知级别（info, warning, error, critical）
	Level string `json:"level"`
	// 通知标题
	Title string `json:"title"`
	// 通知内容
	Content string `json:"content"`
	// 通知时间
	NotificationTime time.Time `json:"notification_time"`
	// 相关服务
	Service string `json:"service,omitempty"`
	// 相关资源
	Resource string `json:"resource,omitempty"`
	// 通知标签
	Tags []string `json:"tags,omitempty"`
}

// NewFileUploadedMessage 创建文件上传完成消息
func NewFileUploadedMessage(fileID, fileName string, fileSize int64, contentType, userID, storagePath string) FileUploadedMessage {
	now := time.Now()
	return FileUploadedMessage{
		BaseMessage: BaseMessage{
			MessageID:   GenerateMessageID(),
			MessageType: "file.uploaded",
			SendTime:    now,
			Source:      "upload-service",
			Version:     "1.0",
		},
		FileID:      fileID,
		FileName:    fileName,
		FileSize:    fileSize,
		ContentType: contentType,
		UserID:      userID,
		UploadTime:  now,
		StoragePath: storagePath,
	}
}

// NewFileProcessedMessage 创建文件处理完成消息
func NewFileProcessedMessage(fileID, processType, result, outputPath string, duration int64, status string) FileProcessedMessage {
	now := time.Now()
	return FileProcessedMessage{
		BaseMessage: BaseMessage{
			MessageID:   GenerateMessageID(),
			MessageType: "file.processed",
			SendTime:    now,
			Source:      "process-service",
			Version:     "1.0",
		},
		FileID:      fileID,
		ProcessType: processType,
		Result:      result,
		ProcessTime: now,
		OutputPath:  outputPath,
		Duration:    duration,
		Status:      status,
	}
}

// NewUserActivityMessage 创建用户活动消息
func NewUserActivityMessage(userID, activityType string, details map[string]interface{}) UserActivityMessage {
	now := time.Now()
	return UserActivityMessage{
		BaseMessage: BaseMessage{
			MessageID:   GenerateMessageID(),
			MessageType: "user.activity",
			SendTime:    now,
			Source:      "user-service",
			Version:     "1.0",
		},
		UserID:       userID,
		ActivityType: activityType,
		ActivityTime: now,
		Details:      details,
	}
}

// NewSystemNotificationMessage 创建系统通知消息
func NewSystemNotificationMessage(level, title, content, service string) SystemNotificationMessage {
	now := time.Now()
	return SystemNotificationMessage{
		BaseMessage: BaseMessage{
			MessageID:   GenerateMessageID(),
			MessageType: "system.notification",
			SendTime:    now,
			Source:      service,
			Version:     "1.0",
		},
		Level:            level,
		Title:            title,
		Content:          content,
		NotificationTime: now,
		Service:          service,
	}
}

// GenerateMessageID 生成消息ID
func GenerateMessageID() string {
	// 简单实现，实际可以使用UUID或其他ID生成算法
	return "msg-" + time.Now().Format("20060102150405.000") + "-" + RandomString(6)
}

// RandomString 生成随机字符串
func RandomString(length int) string {
	// 简单实现，实际可以使用更复杂的随机算法
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // 确保每个字符都不同
	}
	return string(result)
}
