package read

import (
	"api/internal/svc"
	"api/internal/types"
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	db         *sql.DB
	logfile    string
	page       int
	match      string
	path       string
	host       string
	clean      string
	posit      string
	totalLines int
	matchLines int
	contents   []string
	totalPages int
	tasks      []func()
	taskErrors []string
	sshClient  *ssh.Client
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get log info
func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.GetLogInfoReq) (resp *types.GetLogInfoRes, err error) {
	// todo: add your logic here and delete this line
	if req.Logfile == "" || req.Path == "" || req.Host == "" {
		//http.Error(w, "Bad GET param", http.StatusBadRequest)
		return
	}

	if req.Match != "" {
		_, err := regexp.Compile(req.Match)
		if err != nil {
			//http.Error(w, "Incorrect regular expression", http.StatusBadRequest)
			return
		}
	}

	if req.Posit != "head" && req.Posit != "tail" {
		//http.Error(w, "Invalid posit", http.StatusBadRequest)
		return
	}

	err = readLogfile()
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func readLogfile() error {
	tasks = []func(){
		makeTotalLines,
		makeMatchLines,
		makeContents,
	}

	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(task func()) {
			defer wg.Done()
			task()
		}(task)
	}
	wg.Wait()

	if len(taskErrors) > 0 {
		return fmt.Errorf(strings.Join(taskErrors, "; "))
	}

	makeTotalPages()
	return nil
}

func makeContents() {
	var readCmd string
	if match != "" && clean == "true" {
		readCmd = fmt.Sprintf("grep -a \"%s\" %s | head -n %d | tail -n +%d", match, path, page*1000, (page-1)*1000+1)
	} else {
		if posit == "head" {
			readCmd = fmt.Sprintf("head -n %d %s | tail -n +%d", page*1000, path, (page-1)*1000+1)
		} else {
			readCmd = fmt.Sprintf("tail -n +%d %s | head -n 1000", ((page-1)*1000)+1, path)
		}
	}

	output, err := command(readCmd)
	if err != nil {
		taskErrors = append(taskErrors, fmt.Sprintf("get contents error, %v", err))
	} else {
		contents = strings.Split(strings.TrimSpace(output), "\n")
	}
}

func makeTotalLines() {
	output, err := command(fmt.Sprintf("wc -l %s", path))
	if err != nil {
		taskErrors = append(taskErrors, fmt.Sprintf("get lines error, %v", err))
	} else {
		totalLines, _ = strconv.Atoi(strings.Fields(output)[0])
	}
}

func makeMatchLines() {
	if page == 1 && match == "" {
		matchLines = totalLines
	} else if page == 1 {
		output, err := command(fmt.Sprintf("grep -c \"%s\" %s", match, path))
		if err != nil {
			taskErrors = append(taskErrors, fmt.Sprintf("get grep lines error, %v", err))
		} else {
			matchLines, _ = strconv.Atoi(strings.Fields(output)[0])
		}
	}
}

func makeTotalPages() {
	if page == 1 && clean == "true" {
		totalPages = (matchLines / 1000) + 1
		if matchLines%1000 == 0 {
			totalPages--
		}
	} else if page == 1 {
		totalPages = (totalLines / 1000) + 1
		if totalLines%1000 == 0 {
			totalPages--
		}
	}
}

func command(cmd string) (string, error) {
	var output []byte
	var err error
	if host == "localhost" || host == "127.0.0.1" {
		output, err = exec.Command("sh", "-c", cmd).CombinedOutput()
	} else {
		session, err := sshClient.NewSession()
		if err != nil {
			return "", err
		}
		defer session.Close()
		output, err = session.CombinedOutput(cmd)
	}
	return string(output), err
}
