package service

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	sshx "github.com/Auroraol/cloud-storage/common/ssh"
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
	//client, err := sshx.NewClient("101.37.165.220:22", credential, sshx.SetEstablishTimeout(10*time.Second), sshx.SetLogger(sshx.DefaultLogger{}))
	//if err != nil {
	//	panic(err)
	//}

	//_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
	now := time.Now()
	startTime := now.Add(-time.Duration(timeRange) * time.Hour)
	startTimeStr := startTime.Format("2006-01-02T15:04:05")

	// 构建命令
	cmd := fmt.Sprintf("awk '$1 >= \"%s\" {print $0}' %s", startTimeStr, path)

	// 执行命令
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	// 日志处理
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 解析时间戳
		var timestamp int64
		parts := strings.SplitN(line, " ", 2)
		if len(parts) > 1 {
			if t, err := time.Parse("2006-01-02T15:04:05", parts[0]); err == nil {
				timestamp = t.Unix()
			} else {
				// 无法解析时间时跳过该行
				continue
			}
		} else {
			continue
		}

		// 统一转为小写进行匹配检查
		lineLower := strings.ToLower(line)

		// 检查每个监控项
		for _, item := range monitorItems {
			switch strings.ToLower(item) {
			case "requests":
				if strings.Contains(strings.ToLower(lineLower), "request") {
					counter[item][timestamp]++
				}
			case "errors":
				if strings.Contains(strings.ToLower(lineLower), "error") {
					counter[item][timestamp]++
				}
			case "response_time":
				if strings.Contains(strings.ToLower(lineLower), "response_time") {
					// 提取数值示例（假设日志格式为 response_time=0.234）
					if val, err := extractResponseTime(line); err == nil {
						counter[item][timestamp] += val
					} else {
						counter[item][timestamp]++
					}
				}
			}
		}

	}

	// 转换为最终结果并排序
	result := make(map[string][]map[string]interface{})
	for item, timeMap := range counter {
		var series []map[string]interface{}

		// 收集所有时间戳
		timestamps := make([]int64, 0, len(timeMap))
		for ts := range timeMap {
			timestamps = append(timestamps, ts)
		}
		sort.Slice(timestamps, func(i, j int) bool { return timestamps[i] < timestamps[j] })

		// 生成有序序列
		for _, ts := range timestamps {
			series = append(series, map[string]interface{}{
				"timestamp": ts,
				"value":     timeMap[ts],
			})
		}

		result[item] = series
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
