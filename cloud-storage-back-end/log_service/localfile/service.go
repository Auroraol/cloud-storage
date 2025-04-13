package localfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalFileService 提供本地文件服务功能
type LocalFileService struct {
	reader *LocalFileReader
}

// NewLocalFileService 创建一个新的本地文件服务
func NewLocalFileService(basePath string) (*LocalFileService, error) {
	reader, err := NewLocalFileReader(basePath)
	if err != nil {
		return nil, fmt.Errorf("创建本地文件读取器失败: %w", err)
	}

	return &LocalFileService{
		reader: reader,
	}, nil
}

// ListFiles 列出指定目录下的所有文件
func (s *LocalFileService) ListFiles(dirPath string) ([]FileInfo, error) {
	return s.reader.ListFiles(dirPath)
}

// GetFileStat 获取目录的统计信息
func (s *LocalFileService) GetFileStat(dirPath string) (*FileStat, error) {
	fullPath := filepath.Join(s.reader.BasePath, dirPath)

	// 验证路径是否存在
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法访问路径: %w", err)
	}

	if !info.IsDir() {
		return nil, errors.New("提供的路径不是一个目录")
	}

	// 统计信息
	stat := &FileStat{
		TotalFiles:     0,
		TotalDirs:      0,
		TotalSize:      0,
		LogFileCount:   0,
		RecentModified: 0,
	}

	// 当前时间，用于计算最近修改的文件
	now := time.Now()
	recentTime := now.Add(-24 * time.Hour) // 24小时内

	// 遍历目录
	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过访问出错的文件
		}

		if info.IsDir() {
			if path != fullPath { // 不计算根目录
				stat.TotalDirs++
			}
		} else {
			stat.TotalFiles++
			stat.TotalSize += info.Size()

			// 检查是否是日志文件
			if strings.HasSuffix(strings.ToLower(info.Name()), ".log") {
				stat.LogFileCount++
			}

			// 检查是否是最近修改的文件
			if info.ModTime().After(recentTime) {
				stat.RecentModified++
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return stat, nil
}

// ReadLogFile 读取日志文件并按条件过滤
func (s *LocalFileService) ReadLogFile(filePath string, filter LogFilter) ([]LogEntry, error) {
	fullPath := filepath.Join(s.reader.BasePath, filePath)

	// 验证路径是否存在
	_, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法访问文件: %w", err)
	}

	// 读取文件内容
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 解析日志
	entries := make([]LogEntry, 0)
	scanner := NewLogScanner(file)

	for scanner.Scan() {
		entry := scanner.Entry()

		// 应用过滤条件
		if s.filterLogEntry(entry, filter) {
			entries = append(entries, entry)

			// 检查是否达到最大结果数
			if filter.MaxResults > 0 && len(entries) >= filter.MaxResults {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取日志失败: %w", err)
	}

	return entries, nil
}

// filterLogEntry 根据过滤条件过滤日志条目
func (s *LocalFileService) filterLogEntry(entry LogEntry, filter LogFilter) bool {
	// 时间范围过滤
	if filter.StartTime != nil && entry.Timestamp.Before(*filter.StartTime) {
		return false
	}

	if filter.EndTime != nil && entry.Timestamp.After(*filter.EndTime) {
		return false
	}

	// 级别过滤
	if filter.Level != "" && !strings.EqualFold(entry.Level, filter.Level) {
		return false
	}

	// 关键词过滤
	if filter.Keyword != "" && !strings.Contains(strings.ToLower(entry.Content), strings.ToLower(filter.Keyword)) {
		return false
	}

	return true
}

// ReadFileContent 读取文件内容
func (s *LocalFileService) ReadFileContent(filePath string, offset, limit int64) (string, error) {
	return s.reader.ReadFileContent(filePath, offset, limit)
}
