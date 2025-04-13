package localfile

import (
	"bufio"
	"io"
	"regexp"
	"time"
)

// LogScanner 日志扫描器
type LogScanner struct {
	scanner *bufio.Scanner
	entry   LogEntry
	lineNum int
	err     error
}

// 常见日志格式的正则表达式
var (
	// 标准格式: 2023-04-01T12:34:56.789Z INFO [Module] - Message
	standardLogRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2}))\s+(\w+)\s+(?:\[([^\]]+)\]\s+)?-?\s*(.*)$`)

	// 简单格式: 2023/04/01 12:34:56 [INFO] Message
	simpleLogRegex = regexp.MustCompile(`^(\d{4}/\d{2}/\d{2}\s+\d{2}:\d{2}:\d{2})\s+\[(\w+)\]\s+(.*)$`)
)

// NewLogScanner 创建一个新的日志扫描器
func NewLogScanner(r io.Reader) *LogScanner {
	return &LogScanner{
		scanner: bufio.NewScanner(r),
		lineNum: 0,
	}
}

// Scan 扫描下一行日志
func (s *LogScanner) Scan() bool {
	// 扫描下一行
	if !s.scanner.Scan() {
		s.err = s.scanner.Err()
		return false
	}

	s.lineNum++
	line := s.scanner.Text()

	// 尝试解析不同的日志格式
	var timestamp time.Time
	var level, content string
	var matched bool

	// 尝试标准格式
	if matches := standardLogRegex.FindStringSubmatch(line); len(matches) >= 5 {
		timestamp, s.err = time.Parse(time.RFC3339, matches[1])
		if s.err != nil {
			// 尝试其他时间格式
			timestamp, s.err = time.Parse("2006-01-02T15:04:05.000Z", matches[1])
			if s.err != nil {
				timestamp = time.Now() // 如果解析失败，使用当前时间
			}
		}
		level = matches[2]
		content = matches[4]
		matched = true
	}

	// 尝试简单格式
	if !matched {
		if matches := simpleLogRegex.FindStringSubmatch(line); len(matches) >= 4 {
			timestamp, s.err = time.Parse("2006/01/02 15:04:05", matches[1])
			if s.err != nil {
				timestamp = time.Now() // 如果解析失败，使用当前时间
			}
			level = matches[2]
			content = matches[3]
			matched = true
		}
	}

	// 如果所有格式都不匹配，将整行视为内容
	if !matched {
		timestamp = time.Now()
		level = "INFO" // 默认级别
		content = line
	}

	// 更新当前条目
	s.entry = LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Content:   content,
		LineNum:   s.lineNum,
	}

	return true
}

// Entry 获取当前解析的日志条目
func (s *LogScanner) Entry() LogEntry {
	return s.entry
}

// Err 获取扫描过程中的错误
func (s *LogScanner) Err() error {
	return s.err
}

// Scanner 简单的行扫描器
type Scanner struct {
	scanner *bufio.Scanner
	err     error
}

// NewScanner 创建一个新的行扫描器
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(r),
	}
}

// Scan 扫描下一行
func (s *Scanner) Scan() bool {
	if !s.scanner.Scan() {
		s.err = s.scanner.Err()
		return false
	}
	return true
}

// Text 获取当前行的文本
func (s *Scanner) Text() string {
	return s.scanner.Text()
}

// Err 获取扫描过程中的错误
func (s *Scanner) Err() error {
	return s.err
}
