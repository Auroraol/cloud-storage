package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveUploadedFile 保存上传的文件到临时目录
func SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "upload-*"+filepath.Ext(file.Filename))
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// 复制文件内容
	_, err = io.Copy(tempFile, src)
	if err != nil {
		os.Remove(tempFile.Name()) // 清理临时文件
		return "", err
	}

	return tempFile.Name(), nil
}

// ParseMetadata 解析元数据JSON字符串
func ParseMetadata(metadataStr string) (map[string]string, error) {
	if metadataStr == "" {
		return nil, nil
	}

	metadata := make(map[string]string)
	err := json.Unmarshal([]byte(metadataStr), &metadata)
	if err != nil {
		return nil, fmt.Errorf("解析元数据失败: %v", err)
	}

	return metadata, nil
}

// CleanupTempFile 清理临时文件
func CleanupTempFile(filePath string) {
	if filePath != "" {
		os.Remove(filePath)
	}
}
