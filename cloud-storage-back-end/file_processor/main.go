package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Auroraol/cloud-storage/common/logx"
	"github.com/Auroraol/cloud-storage/common/mq"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	configFile = flag.String("f", ".env", "配置文件路径")
)

func main() {
	flag.Parse()

	// 加载环境变量
	if err := godotenv.Load(*configFile); err != nil {
		fmt.Printf("无法加载配置文件: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	logConfig := logx.LogConfig{
		LogLevel:          "info",
		LogFormat:         "json",
		LogPath:           "./logs",
		LogFileName:       "file_processor.log",
		LogFileMaxSize:    100,
		LogFileMaxBackups: 10,
		LogMaxAge:         30,
		LogCompress:       true,
		LogStdout:         true,
		SeparateLevel:     true,
	}
	if err := logx.InitLogger(logConfig); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	// 获取Pulsar配置
	pulsarURL := os.Getenv("PULSAR_URL")
	if pulsarURL == "" {
		pulsarURL = "pulsar://localhost:6650" // 默认值
	}

	// 创建Pulsar管理器
	pulsarManager, err := mq.NewPulsarManager(mq.PulsarConfig{
		URL: pulsarURL,
	})
	if err != nil {
		zap.S().Fatalf("创建Pulsar管理器失败: %v", err)
	}
	defer pulsarManager.Close()

	zap.S().Info("文件处理服务启动成功")

	// 启动文件上传消息消费者
	err = mq.StartFileUploadedConsumer(pulsarManager, handleFileUploaded)
	if err != nil {
		zap.S().Fatalf("启动文件上传消息消费者失败: %v", err)
	}

	zap.S().Info("文件上传消息消费者启动成功")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.S().Info("正在关闭文件处理服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 在这里可以添加优雅关闭的逻辑
	<-ctx.Done()

	zap.S().Info("文件处理服务已关闭")
}

// 处理文件上传消息
func handleFileUploaded(msg mq.FileUploadedMessage) error {
	zap.S().Infof("收到文件上传消息: 文件ID=%s, 文件名=%s, 大小=%d, 用户ID=%s",
		msg.FileID, msg.FileName, msg.FileSize, msg.UserID)

	// 这里可以添加文件处理逻辑，例如：
	// 1. 生成缩略图
	// 2. 提取文件元数据
	// 3. 文件格式转换
	// 4. 内容分析
	// 5. 更新搜索索引
	// 等等

	// 模拟处理时间
	time.Sleep(500 * time.Millisecond)

	zap.S().Infof("文件处理完成: 文件ID=%s", msg.FileID)
	return nil
}
