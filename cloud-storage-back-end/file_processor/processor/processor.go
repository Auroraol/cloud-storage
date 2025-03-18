package processor

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Auroraol/cloud-storage/common/mq/pulsar"
	"go.uber.org/zap"
)

// FileProcessor 文件处理器接口
type FileProcessor interface {
	Process(msg pulsar.FileUploadedMessage) error
}

// 根据文件类型获取处理器
func GetProcessorForFile(fileName string) FileProcessor {
	ext := strings.ToLower(filepath.Ext(fileName))

	switch {
	case isImageExt(ext):
		return &ImageProcessor{}
	case isVideoExt(ext):
		return &VideoProcessor{}
	case isDocumentExt(ext):
		return &DocumentProcessor{}
	default:
		return &DefaultProcessor{}
	}
}

// 判断文件类型的辅助函数
func isImageExt(ext string) bool {
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff"}
	return contains(imageExts, ext)
}

func isVideoExt(ext string) bool {
	videoExts := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv", ".webm"}
	return contains(videoExts, ext)
}

func isDocumentExt(ext string) bool {
	docExts := []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt"}
	return contains(docExts, ext)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ImageProcessor 图片处理器
type ImageProcessor struct{}

func (p *ImageProcessor) Process(msg pulsar.FileUploadedMessage) error {
	zap.S().Infof("开始处理图片文件: %s", msg.FileName)

	// 模拟处理时间
	time.Sleep(500 * time.Millisecond)

	// 这里实现图片处理逻辑，例如：
	// 1. 下载原始图片
	// 2. 生成缩略图
	// 3. 提取图片元数据
	// 4. 上传处理后的图片
	// 5. 更新数据库记录

	zap.S().Infof("图片处理完成: %s", msg.FileName)
	return nil
}

// VideoProcessor 视频处理器
type VideoProcessor struct{}

func (p *VideoProcessor) Process(msg pulsar.FileUploadedMessage) error {
	zap.S().Infof("开始处理视频文件: %s", msg.FileName)

	// 模拟处理时间
	time.Sleep(1 * time.Second)

	// 这里实现视频处理逻辑，例如：
	// 1. 下载原始视频
	// 2. 生成视频预览图
	// 3. 提取视频元数据（时长、分辨率等）
	// 4. 可能的转码操作
	// 5. 上传处理后的文件
	// 6. 更新数据库记录

	zap.S().Infof("视频处理完成: %s", msg.FileName)
	return nil
}

// DocumentProcessor 文档处理器
type DocumentProcessor struct{}

func (p *DocumentProcessor) Process(msg pulsar.FileUploadedMessage) error {
	zap.S().Infof("开始处理文档文件: %s", msg.FileName)

	// 模拟处理时间
	time.Sleep(300 * time.Millisecond)

	// 这里实现文档处理逻辑，例如：
	// 1. 下载原始文档
	// 2. 生成预览（如PDF预览）
	// 3. 提取文本内容用于搜索
	// 4. 上传处理后的文件
	// 5. 更新数据库记录

	zap.S().Infof("文档处理完成: %s", msg.FileName)
	return nil
}

// DefaultProcessor 默认处理器
type DefaultProcessor struct{}

func (p *DefaultProcessor) Process(msg pulsar.FileUploadedMessage) error {
	zap.S().Infof("处理其他类型文件: %s", msg.FileName)

	// 对于未识别的文件类型，可能只需要记录基本信息
	// 或者执行一些通用的处理

	zap.S().Infof("文件处理完成: %s", msg.FileName)
	return nil
}

// ProcessFile 处理文件的主函数
func ProcessFile(msg pulsar.FileUploadedMessage) error {
	// 获取适合的处理器
	processor := GetProcessorForFile(msg.FileName)

	// 记录处理开始
	startTime := time.Now()
	zap.S().Infof("开始处理文件: ID=%s, 文件名=%s", msg.FileID, msg.FileName)

	// 执行处理
	err := processor.Process(msg)

	// 记录处理结束
	duration := time.Since(startTime)
	if err != nil {
		zap.S().Errorf("文件处理失败: ID=%s, 错误=%v, 耗时=%v", msg.FileID, err, duration)
		return fmt.Errorf("处理文件失败: %w", err)
	}

	zap.S().Infof("文件处理成功: ID=%s, 耗时=%v", msg.FileID, duration)
	return nil
}
