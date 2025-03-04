package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	stdtime "time"

	sshx "github.com/Auroraol/cloud-storage/common/ssh"
	"github.com/Auroraol/cloud-storage/common/time"
)

// SSHService SSH服务接口
type SSHService interface {
	// 连接主机
	Connect(host, port, user, password, privateKeyPath string) error
	// 读取日志文件
	ReadLogFile(path string, match string, page, pageSize int) ([]string, int, error)
	// 获取日志文件列表
	GetLogFiles(path string) ([]string, error)
	// 监控日志文件
	MonitorLogFile(path string, monitorItems []string, timeRange int) (map[string][]map[string]interface{}, error)
	// 关闭连接
	Close()
}

// sshService SSH服务实现
type sshService struct {
	client sshx.Client
	host   string
}

// NewSSHService 创建SSH服务
func NewSSHService() SSHService {
	return &sshService{}
}

// Connect 连接主机
func (s *sshService) Connect(host, port, user, password, privateKeyPath string) error {
	// 创建凭证
	credential := sshx.Credential{
		User:           user,
		Password:       password,
		PrivateKeyPath: privateKeyPath,
	}

	// 创建SSH客户端
	client, err := sshx.NewClient(host+":"+port, credential)
	if err != nil {
		return fmt.Errorf("连接主机失败: %v", err)
	}

	//privateKeyConf := sshx.Credential{User: "root", Password: "-+66..[]l"}
	//
	//// 创建 sshx 客户端
	//client, err := sshx.NewClient("101.37.165.220:22", credential, sshx.SetEstablishTimeout(10*stdtime.Second), sshx.SetLogger(sshx.DefaultLogger{}))
	//if err != nil {
	//	panic(err)
	//}

	//_, cancel := context.WithTimeout(context.Background(), 5*stdtime.Second)
	//defer cancel()
	if err := client.Handle(func(sub sshx.EnhanceClient) error {
		// if _, err := sub.ReceiveFile("/tmp/xxx", "/etc/passwd", false, true); err != nil {
		// 	panic(err)
		// }

		path := "/opt/goTest/log"

		data, err := sub.ReadFile(path)
		if err != nil {
			panic(err)
		}

		fmt.Printf("filebeat yml: %s", string(data))
		return nil
	}); err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}
	s.client = client
	s.host = host
	return nil
}

// ReadLogFile 读取日志文件
func (s *sshService) ReadLogFile(path string, match string, page, pageSize int) ([]string, int, error) {
	if s.client == nil {
		return nil, 0, fmt.Errorf("未连接主机")
	}

	// 构建命令
	cmd := fmt.Sprintf("cat %s", path)
	if match != "" {
		cmd = fmt.Sprintf("grep -i '%s' %s", match, path)
	}

	// 执行命令
	ctx, cancel := context.WithTimeout(context.Background(), 30*stdtime.Second)
	defer cancel()

	output, err := s.client.Command(ctx, cmd)
	if err != nil {
		return nil, 0, fmt.Errorf("读取日志文件失败: %v", err)
	}

	// 分割日志行
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// 计算总行数
	totalLines := len(lines)

	// 分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= totalLines {
		return []string{}, totalLines, nil
	}
	if end > totalLines {
		end = totalLines
	}

	return lines[start:end], totalLines, nil
}

// GetLogFiles 获取日志文件列表
func (s *sshService) GetLogFiles(path string) ([]string, error) {
	if s.client == nil {
		return nil, fmt.Errorf("未连接主机")
	}

	// 构建命令
	cmd := fmt.Sprintf("find %s -type f -name '*.log' | sort", path)

	// 执行命令
	ctx, cancel := context.WithTimeout(context.Background(), 30*stdtime.Second)
	defer cancel()

	output, err := s.client.Command(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("获取日志文件列表失败: %v", err)
	}

	// 分割文件列表
	files := strings.Split(string(output), "\n")
	if len(files) > 0 && files[len(files)-1] == "" {
		files = files[:len(files)-1]
	}

	return files, nil
}

// MonitorLogFile 监控日志文件
func (s *sshService) MonitorLogFile(path string, monitorItems []string, timeRange int) (map[string][]map[string]interface{}, error) {
	if s.client == nil {
		return nil, fmt.Errorf("未连接主机")
	}

	// 计算时间范围
	now := time.LocalTimeNow()
	startTime := now.Add(-stdtime.Duration(timeRange) * stdtime.Hour)
	startTimeStr := startTime.Format("2006-01-02 15:04:05")

	// 构建命令
	cmd := fmt.Sprintf("awk '$1 >= \"%s\" {print $0}' %s", startTimeStr, path)

	// 执行命令
	ctx, cancel := context.WithTimeout(context.Background(), 30*stdtime.Second)
	defer cancel()

	output, err := s.client.Command(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("监控日志文件失败: %v", err)
	}

	// 使用双层 map 进行计数 [监控项][时间戳]计数
	counter := make(map[string]map[int64]int)
	for _, item := range monitorItems {
		counter[item] = make(map[int64]int)
	}

	// 创建一个map来存储调用者统计信息
	callerStats := make(map[string]int)

	// 日志处理
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		// 尝试解析为JSON格式
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err == nil {
			// 提取时间
			var logTime stdtime.Time
			var toTimestamp int64
			if timeStr, ok := logEntry["time"].(string); ok {
				// 尝试解析时间，支持多种格式
				toTimestamp, err = time.StringTimeToTimestamp(timeStr)
				if err != nil {
					// 如果解析失败，使用当前时间
					logTime = time.LocalTimeNow()
				} else {
					logTime = stdtime.Unix(toTimestamp, 0)
				}
			} else {
				// 如果没有时间字段，使用当前时间
				logTime = time.LocalTimeNow()
			}

			// 检查日志是否在指定的时间范围内
			if logTime.Before(startTime) {
				continue // 跳过不在时间范围内的日志
			}

			// 提取日志级别
			var level string
			if levelValue, ok := logEntry["level"].(string); ok {
				level = strings.ToUpper(levelValue)
			}

			// 提取消息
			var message string
			if msgValue, ok := logEntry["msg"].(string); ok {
				message = msgValue
			}

			// 提取调用者信息并统计
			if callerValue, ok := logEntry["caller"].(string); ok {
				// 统计调用者
				callerStats[callerValue]++
			}

			// 按分钟取整
			timestamp := toTimestamp

			// 按照监控项进行统计
			for _, item := range monitorItems {
				switch item {
				case "requests":
					// 统计所有请求
					counter["requests"][timestamp]++
				case "errors":
					// 统计错误日志
					if level == "ERROR" {
						counter["errors"][timestamp]++
					}
				case "response_time":
					// 尝试从消息中提取响应时间
					if responseTime, err := extractResponseTime(message); err == nil {
						counter["response_time"][timestamp] += responseTime
					}
				case "debug_logs":
					// 统计调试日志
					if level == "DEBUG" {
						counter["debug_logs"][timestamp]++
					}
				case "warn_logs":
					// 统计警告日志
					if level == "WARN" {
						counter["warn_logs"][timestamp]++
					}
				case "info_logs":
					// 统计信息日志
					if level == "INFO" {
						counter["info_logs"][timestamp]++
					}
				}
			}
		}
	}

	// 构建结果
	result := make(map[string][]map[string]interface{})
	for item, timestamps := range counter {
		data := make([]map[string]interface{}, 0, len(timestamps))
		for ts, count := range timestamps {
			data = append(data, map[string]interface{}{
				"timestamp": ts,
				"value":     count,
			})
		}

		// 按时间戳排序
		sort.Slice(data, func(i, j int) bool {
			return data[i]["timestamp"].(int64) < data[j]["timestamp"].(int64)
		})

		result[item] = data
	}

	// 添加调用者统计信息
	if len(callerStats) > 0 {
		callerData := make([]map[string]interface{}, 0, len(callerStats))
		for caller, count := range callerStats {
			callerData = append(callerData, map[string]interface{}{
				"caller": caller,
				"value":  count,
			})
		}
		result["caller_stats"] = callerData
	}

	return result, nil
}

// 辅助函数：从日志行提取响应时间（示例实现）
func extractResponseTime(line string) (int, error) {
	// 假设日志格式包含 response_time=123ms
	re := regexp.MustCompile(`response_time=(\d+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 2 {
		return 0, fmt.Errorf("未找到响应时间")
	}

	val, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return val, nil
}

// Close 关闭连接
func (s *sshService) Close() {
	// SSH客户端没有显式的关闭方法，连接会在使用后自动关闭
	s.client = nil
}
