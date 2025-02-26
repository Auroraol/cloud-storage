package logic

//
//import (
//	"context"
//	"fmt"
//	"golang.org/x/crypto/ssh"
//	"os/exec"
//	"regexp"
//	"strconv"
//	"strings"
//	"sync"
//
//	"cloud-storage/log_service/api/internal/svc"
//	"cloud-storage/log_service/api/internal/types"
//
//	"github.com/zeromicro/go-zero/core/logx"
//)
//
//type ReadLogic struct {
//	logx.Logger
//	ctx    context.Context
//	svcCtx *svc.ServiceContext
//}
//
//// 日志阅读
//func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
//	return &ReadLogic{
//		Logger: logx.WithContext(ctx),
//		ctx:    ctx,
//		svcCtx: svcCtx,
//	}
//}
//
//func (l *ReadLogic) Read(req *types.GetLogInfoReq) (resp *types.GetLogInfoRes, err error) {
//	// 参数验证
//	if req.Logfile == "" {
//		return nil, fmt.Errorf("logfile is required")
//	}
//	if req.Path == "" {
//		return nil, fmt.Errorf("path is required")
//	}
//	if req.Host == "" {
//		return nil, fmt.Errorf("host is required")
//	}
//	if req.Match != "" {
//		if _, err := regexp.Compile(req.Match); err != nil {
//			return nil, fmt.Errorf("incorrect regular expression")
//		}
//	}
//	if req.Page <= 0 {
//		return nil, fmt.Errorf("invalid page")
//	}
//	if req.Posit != "head" && req.Posit != "tail" {
//		return nil, fmt.Errorf("invalid posit")
//	}
//
//	logfileRow, err := l.getLogfileInfo(req.Logfile)
//	if err != nil {
//		return nil, err
//	}
//	if logfileRow == nil {
//		return nil, fmt.Errorf("logfile not exist")
//	}
//	if !strings.Contains(logfileRow.Host, req.Host) {
//		return nil, fmt.Errorf("invalid host")
//	}
//
//	// SSH 连接
//	var sshClient *ssh.Client
//	if req.Host != "localhost" && req.Host != "127.0.0.1" {
//		sshClient, err = l.connectSSH(req.Host)
//		if err != nil {
//			return nil, err
//		}
//		defer sshClient.Close()
//	}
//
//	// 并发执行命令
//	var wg sync.WaitGroup
//	var totalLines, matchLines int
//	var contents []string
//	var totalLinesErr, matchLinesErr, contentsErr error
//
//	wg.Add(3)
//	go func() {
//		defer wg.Done()
//		totalLines, totalLinesErr = l.makeTotalLines(req.Path, sshClient)
//	}()
//	go func() {
//		defer wg.Done()
//		matchLines, matchLinesErr = l.makeMatchLines(req.Path, req.Match, sshClient)
//	}()
//	go func() {
//		defer wg.Done()
//		contents, contentsErr = l.makeContents(req.Path, req.Match, req.Clean, req.Page, req.Posit, sshClient)
//	}()
//	wg.Wait()
//
//	if totalLinesErr != nil {
//		return nil, totalLinesErr
//	}
//	if matchLinesErr != nil {
//		return nil, matchLinesErr
//	}
//	if contentsErr != nil {
//		return nil, contentsErr
//	}
//
//	totalPages := (totalLines / 1000) + 1
//	if totalLines%1000 == 0 {
//		totalPages--
//	}
//
//	resp = &types.GetLogInfoRes{
//		Contents:   contents,
//		Page:       req.Page,
//		TotalPages: totalPages,
//		TotalLines: totalLines,
//		MatchLines: matchLines,
//		Lines:      len(contents),
//	}
//	return resp, nil
//}
//
//func (l *ReadLogic) getLogfileInfo(logfile string) (*logfileRow, error) {
//	// 这里需要实现从数据库或其他存储中获取 logfile 信息的逻辑
//	// 示例返回一个固定的 logfileRow
//	return &logfileRow{
//		Host: "localhost",
//	}, nil
//}
//
//type logfileRow struct {
//	Host string
//}
//
//func (l *ReadLogic) connectSSH(host string) (*ssh.Client, error) {
//	config := &ssh.ClientConfig{
//		User: "your_user",
//		Auth: []ssh.AuthMethod{
//			ssh.Password("your_password"),
//		},
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//	}
//	client, err := ssh.Dial("tcp", host+":22", config)
//	if err != nil {
//		return nil, err
//	}
//	return client, nil
//}
//
//func (l *ReadLogic) makeTotalLines(path string, sshClient *ssh.Client) (int, error) {
//	cmd := "wc -l " + path
//	output, err := l.executeCommand(cmd, sshClient)
//	if err != nil {
//		return 0, err
//	}
//	lines := strings.Split(output, " ")
//	totalLines, err := strconv.Atoi(strings.TrimSpace(lines[0]))
//	if err != nil {
//		return 0, err
//	}
//	return totalLines, nil
//}
//
//func (l *ReadLogic) makeMatchLines(path, match string, sshClient *ssh.Client) (int, error) {
//	cmd := fmt.Sprintf("grep -c \"%s\" %s", match, path)
//	output, err := l.executeCommand(cmd, sshClient)
//	if err != nil {
//		return 0, err
//	}
//	lines := strings.Split(output, " ")
//	matchLines, err := strconv.Atoi(strings.TrimSpace(lines[0]))
//	if err != nil {
//		return 0, err
//	}
//	return matchLines, nil
//}
//
//func (l *ReadLogic) makeContents(path, match, clean string, page int, posit string, sshClient *ssh.Client) ([]string, error) {
//	var cmd string
//	if match != "" && clean == "true" {
//		cmd = fmt.Sprintf("grep -a \"%s\" %s | head -n %d | tail -n +%d", match, path, page*1000, (page-1)*1000+1)
//	} else {
//		if posit == "head" {
//			cmd = fmt.Sprintf("head -n %d %s | tail -n +%d", page*1000, path, (page-1)*1000+1)
//		} else {
//			cmd = fmt.Sprintf("tail -n +%d %s | head -n 1000", ((page-1)*1000)+1, path)
//		}
//	}
//	output, err := l.executeCommand(cmd, sshClient)
//	if err != nil {
//		return nil, err
//	}
//	return strings.Split(output, "\n"), nil
//}
//
//func (l *ReadLogic) executeCommand(cmd string, sshClient *ssh.Client) (string, error) {
//	if sshClient == nil {
//		out, err := exec.Command("sh", "-c", cmd).Output()
//		if err != nil {
//			return "", err
//		}
//		return string(out), nil
//	}
//	session, err := sshClient.NewSession()
//	if err != nil {
//		return "", err
//	}
//	defer session.Close()
//	output, err := session.CombinedOutput(cmd)
//	if err != nil {
//		return "", err
//	}
//	return string(output), nil
//}
