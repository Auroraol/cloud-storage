package oss

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// 普通文件上传的大小限制（20MB）
const NormalUploadLimit = 20 * 1024 * 1024

// 分片大小（1MB）
const PartSize = 1024 * 1024

// FileUploadOptions 文件上传选项
type FileUploadOptions struct {
	FilePath    string            // 本地文件路径
	ObjectKey   string            // OSS对象键
	ContentType string            // 内容类型
	Metadata    map[string]string // 元数据
}

// 普通上传
func NormalUpload(options FileUploadOptions) (string, error) {
	file, err := os.Open(options.FilePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 生成元数据选项
	metaOptions := GenFileMeta(options.Metadata)
	if options.ContentType != "" {
		metaOptions = append(metaOptions, oss.ContentType(options.ContentType))
	}

	bucket := Bucket()
	if bucket == nil {
		return "", fmt.Errorf("获取Bucket失败")
	}

	err = bucket.PutObject(options.ObjectKey, file, metaOptions...)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %v", err)
	}

	return fmt.Sprintf("https://%s.%s/%s", Config.BucketName, Config.Endpoint, options.ObjectKey), nil
}

// 分片上传
func MultipartUpload(options FileUploadOptions) (string, error) {
	bucket := Bucket()
	if bucket == nil {
		return "", fmt.Errorf("获取Bucket失败")
	}

	// 初始化分片上传
	imur, err := bucket.InitiateMultipartUpload(options.ObjectKey)
	if err != nil {
		return "", fmt.Errorf("初始化分片上传失败: %v", err)
	}

	file, err := os.Open(options.FilePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 计算分片数量
	chunks := int(fileInfo.Size() / PartSize)
	if fileInfo.Size()%PartSize > 0 {
		chunks++
	}

	// 创建一个分片上传的切片数组
	parts := make([]oss.UploadPart, 0, chunks)

	// 分片上传
	for i := 1; i <= chunks; i++ {
		// 跳转到每个分片的开头
		file.Seek(PartSize*int64(i-1), 0)
		// 计算每个分片的大小
		partSize := PartSize
		if i == chunks {
			partSize = int(fileInfo.Size() - PartSize*int64(i-1))
		}

		// 上传分片
		part, err := bucket.UploadPart(imur, file, int64(partSize), i)
		if err != nil {
			// 如果上传失败，取消分片上传
			bucket.AbortMultipartUpload(imur)
			return "", fmt.Errorf("上传分片 %d 失败: %v", i, err)
		}
		parts = append(parts, part)
	}

	// 完成分片上传
	_, err = bucket.CompleteMultipartUpload(imur, parts)
	if err != nil {
		return "", fmt.Errorf("完成分片上传失败: %v", err)
	}

	return fmt.Sprintf("https://%s.%s/%s", Config.BucketName, Config.Endpoint, options.ObjectKey), nil
}

// 断点续传上传
func ResumeUpload(options FileUploadOptions) (string, error) {
	bucket := Bucket()
	if bucket == nil {
		return "", fmt.Errorf("获取Bucket失败")
	}

	// 创建断点续传文件
	checkpoint := filepath.Join(os.TempDir(), fmt.Sprintf("%s.cp", filepath.Base(options.FilePath)))

	// 断点续传选项
	uploadOptions := []oss.Option{
		oss.Checkpoint(true, checkpoint), // 开启断点续传
	}

	if options.ContentType != "" {
		uploadOptions = append(uploadOptions, oss.ContentType(options.ContentType))
	}

	// 添加元数据
	if options.Metadata != nil {
		for k, v := range options.Metadata {
			uploadOptions = append(uploadOptions, oss.Meta(k, v))
		}
	}

	// 执行断点续传
	err := bucket.UploadFile(options.ObjectKey, options.FilePath, PartSize, uploadOptions...)
	if err != nil {
		return "", fmt.Errorf("断点续传上传失败: %v", err)
	}

	// 清理断点续传文件
	os.Remove(checkpoint)

	return fmt.Sprintf("https://%s.%s/%s", Config.BucketName, Config.Endpoint, options.ObjectKey), nil
}

// 根据文件大小自动选择上传方式
func UploadFile(options FileUploadOptions) (string, error) {
	fileInfo, err := os.Stat(options.FilePath)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 根据文件大小选择上传方式
	if fileInfo.Size() <= NormalUploadLimit {
		// 小于20MB使用普通上传
		return NormalUpload(options)
	} else {
		// 大于20MB使用断点续传
		return ResumeUpload(options)
	}
}
