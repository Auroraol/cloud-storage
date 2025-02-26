package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	sshx "github.com/Auroraol/cloud-storage/common/ssh"
)

// SSHService SSH服务接口
type SSHService interface {
	// 连接主机
	Connect(host, user, password, privateKeyPath string) error
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
func (s *sshService) Connect(host, user, password, privateKeyPath string) error {
	// 创建凭证
	credential := sshx.Credential{
		User:           user,
		Password:       password,
		PrivateKeyPath: privateKeyPath,
	}

	// 创建SSH客户端
	client, err := sshx.NewClient(host, credential)
	if err != nil {
		return fmt.Errorf("连接主机失败: %v", err)
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

	// 计算开始时间
	now := time.Now()
	startTime := now.Add(-time.Duration(timeRange) * time.Hour)
	startTimeStr := startTime.Format("2006-01-02T15:04:05")

	// 构建命令，使用awk提取时间戳和相关信息
	cmd := fmt.Sprintf("awk '$1 >= \"%s\" {print $0}' %s", startTimeStr, path)

	// 执行命令
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err := s.client.Command(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("监控日志文件失败: %v", err)
	}

	// 分析日志数据
	result := make(map[string][]map[string]interface{})
	lines := strings.Split(string(output), "\n")

	// 初始化结果
	for _, item := range monitorItems {
		result[item] = []map[string]interface{}{}
	}

	// 分析日志行，提取监控项数据
	// 这里使用简化的逻辑，实际应根据日志格式进行解析
	for _, line := range lines {
		if line == "" {
			continue
		}

		// 解析时间戳
		timestamp := time.Now().Unix() // 默认使用当前时间
		parts := strings.SplitN(line, " ", 2)
		if len(parts) > 1 {
			if t, err := time.Parse("2006-01-02T15:04:05", parts[0]); err == nil {
				timestamp = t.Unix()
			}
		}

		// 检查监控项
		for _, item := range monitorItems {
			if item == "请求数" && strings.Contains(line, "request") {
				result[item] = append(result[item], map[string]interface{}{
					"timestamp": timestamp,
					"value":     1,
				})
			} else if item == "错误数" && (strings.Contains(line, "error") || strings.Contains(line, "ERROR")) {
				result[item] = append(result[item], map[string]interface{}{
					"timestamp": timestamp,
					"value":     1,
				})
			} else if item == "响应时间" && strings.Contains(line, "response_time") {
				// 尝试提取响应时间
				// 这里使用简化逻辑，实际应根据日志格式提取
				result[item] = append(result[item], map[string]interface{}{
					"timestamp": timestamp,
					"value":     200, // 默认值
				})
			}
		}
	}

	return result, nil
}

// Close 关闭连接
func (s *sshService) Close() {
	// SSH客户端没有显式的关闭方法，连接会在使用后自动关闭
	s.client = nil
}
