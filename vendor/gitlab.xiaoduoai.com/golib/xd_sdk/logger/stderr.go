// +build !windows

package logger

import (
	"os"
	"syscall"

	"github.com/pkg/errors"
)

// redirectStderr to the file passed in
func redirectStderr(f *os.File) error {
	err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return errors.Wrapf(err, "failed to redirect stderr to file(%s)", f.Name())
	}
	return nil
}
