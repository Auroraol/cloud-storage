package logger

import (
	"log"

	"github.com/pkg/errors"
)

type Standard Options

var _ = ResetStandard()

func ResetStandard(opts ...Option) (err error) {
	return ResetStandardWithOptions(newOptions(opts...))
}

func ResetStandardWithOptions(options Options) (err error) {
	l := StandardLogger()
	if err = initLoggerWithOptions(l, options); err != nil {
		return errors.Wrap(err, "failed to initialize logger")
	}
	return nil
}

// SetStandardLogger 设置默认 Logger 为 l
func SetStandardLogger(l XdLogger, opts ...Option) {
	xdStdLogger = l
	if err := initLoggerWithOptions(xdStdLogger, newOptions(opts...)); err != nil {
		log.Println("SetStandardLogger initLoggerWithOptions fail: ", err)
	}
}

func SetStandardLoggerWithOptions(l XdLogger, options Options) {
	xdStdLogger = l
	if err := initLoggerWithOptions(xdStdLogger, options); err != nil {
		log.Println("SetStandardLoggerWithOptions initLoggerWithOptions fail: ", err)
	}
}
