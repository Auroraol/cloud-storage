//go:build arm64
// +build arm64

package logger

import (
	"os"
	"syscall"

	"github.com/pkg/errors"
)

func redirectStderr(f *os.File) error {
	err := syscall.Dup3(int(f.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		// oldfd == newfd这时可以不处理
		if err == syscall.EINVAL {
			return nil
		}
		return errors.Wrapf(err, "failed to redirect stderr to file(%s)", f.Name())
	}
	return nil
}
