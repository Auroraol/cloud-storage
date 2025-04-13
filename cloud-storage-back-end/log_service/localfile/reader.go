package localfile

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalFileReader 提供本地文件读取功能
type LocalFileReader struct {
	BasePath string // 基础路径
}

// NewLocalFileReader 创建一个新的本地文件读取器
func NewLocalFileReader(basePath string) (*LocalFileReader, error) {
	// 验证路径是否存在
	info, err := os.Stat(basePath)
	if err != nil {
		return nil, fmt.Errorf("无法访问基础路径: %w", err)
	}

	if !info.IsDir() {
		return nil, errors.New("提供的路径不是一个目录")
	}

	return &LocalFileReader{
		BasePath: basePath,
	}, nil
}

// ListFiles 列出指定目录下的所有文件
func (r *LocalFileReader) ListFiles(dirPath string) ([]FileInfo, error) {
	fullPath := filepath.Join(r.BasePath, dirPath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		entryPath := filepath.Join(fullPath, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue // 跳过无法获取信息的文件
		}

		files = append(files, FileInfo{
			Path:      entryPath,
			Name:      entry.Name(),
			Size:      info.Size(),
			IsDir:     entry.IsDir(),
			ModTime:   info.ModTime(),
			Extension: filepath.Ext(entry.Name()),
		})
	}

	return files, nil
}

// ReadFileContent 读取文件内容
func (r *LocalFileReader) ReadFileContent(filePath string, offset, limit int64) (string, error) {
	fullPath := filepath.Join(r.BasePath, filePath)

	file, err := os.Open(fullPath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}

	fileSize := fileInfo.Size()

	// 验证偏移量和限制
	if offset < 0 {
		offset = 0
	}

	if offset >= fileSize {
		return "", nil // 偏移量超出文件大小，返回空内容
	}

	if limit <= 0 || limit > fileSize-offset {
		limit = fileSize - offset
	}

	// 设置读取位置
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("设置文件偏移量失败: %w", err)
	}

	// 读取指定长度的内容
	buffer := make([]byte, limit)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("读取文件内容失败: %w", err)
	}

	return string(buffer[:n]), nil
}
