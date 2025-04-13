package localfile

import (
	"time"
)

// LogEntry 表示一个日志条目
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Content   string    `json:"content"`
	Source    string    `json:"source"` // 文件路径
	LineNum   int       `json:"line_num"`
}

// LogFilter 用于筛选日志的条件
type LogFilter struct {
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	Level      string     `json:"level,omitempty"`
	Keyword    string     `json:"keyword,omitempty"`
	MaxResults int        `json:"max_results,omitempty"`
}

// FileInfo 文件信息
type FileInfo struct {
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	IsDir     bool      `json:"is_dir"`
	ModTime   time.Time `json:"mod_time"`
	Extension string    `json:"extension"`
}

// FileStat 文件统计信息
type FileStat struct {
	TotalFiles     int   `json:"total_files"`
	TotalDirs      int   `json:"total_dirs"`
	TotalSize      int64 `json:"total_size"`
	LogFileCount   int   `json:"log_file_count"`
	RecentModified int   `json:"recent_modified"` // 最近24小时内修改的文件数
}
