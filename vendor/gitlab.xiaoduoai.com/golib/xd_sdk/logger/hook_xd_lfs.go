package logger

import (
	"io"

	"gitlab.xiaoduoai.com/golib/xd_sdk/metadata"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type XdLfsHook struct {
	*lfshook.LfsHook
	isCopyTestLog bool
	writer        io.Writer
}

// 将等级为error(及以上)的日志复制一份写到errWriter。
func NewErrWriterHook(errWriter io.Writer) *XdLfsHook {
	lfsh := NewXdLfsHook(
		lfshook.WriterMap{
			ErrorLevel: errWriter,
			FatalLevel: errWriter,
			PanicLevel: errWriter,
		}, newJSONFormatter())
	lfsh.SetIsCopyTestLog(false) // 压测流量不复制
	lfsh.writer = errWriter
	return lfsh
}

func NewXdLfsHook(output interface{}, formatter logrus.Formatter) *XdLfsHook {
	return &XdLfsHook{
		LfsHook:       lfshook.NewHook(output, formatter),
		isCopyTestLog: false,
	}
}

// 覆盖LfsHook的同名方法，控制压测日志的输出
func (hook *XdLfsHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil && metadata.IsTestFlow(entry.Context) && hook.isCopyTestLog == false { // 压测流量，结束写日志过程
		return nil
	}
	return hook.LfsHook.Fire(entry)
}

func (hook *XdLfsHook) SetIsCopyTestLog(isCopyTestLog bool) {
	hook.isCopyTestLog = isCopyTestLog
}
