package xd_error

import (
	"encoding/json"
	"errors"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

// IXDError 是对 XDError 结构体的抽象，主要用于 http 框架统一处理HTTP接口错误（错误状态和错误码）
// 后期xd_error包反馈的error对象应该逐步替换为接口，不应该直接返回结构体
type IXDError interface {
	error
	GetCode() int
	Unwrap() error
	HTTPStatus() int
}

type HTTPError interface {
	error
	HTTPStatus() int
}

const (
	maxStackDepth = 32
	skipStack     = 2
)

// XDError 主要用于 http 框架统一处理HTTP接口错误（错误状态和错误码）
type XDError struct {
	Code     ErrorCodeType `json:"code"`
	Message  string        `json:"message"`
	RTStacks string        `json:"rt_stacks"`
	err      error
	httpCode int
}

func (e *XDError) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *XDError) GetCode() int {
	return int(e.Code)
}

func (e *XDError) Unwrap() error {
	return e.err
}

// HTTPStatus 返回HTTP状态码
func (e *XDError) HTTPStatus() (code int) {
	if e.httpCode > 0 {
		return e.httpCode
	}
	// 从特定 code 解析 http 状态码。
	// 因可能存在的兼容性问题，暂时注释
	// switch e.Code {
	// case ErrorCodeWrongParam:
	// 	return 400
	// case ErrorCodeSystem:
	// 	return 500
	// }
	return 200
}

func getStacks(skip int) string {
	builder := strings.Builder{}
	pc := make([]uintptr, maxStackDepth)
	n := runtime.Callers(skip, pc)
	for i := 0; i < n-2; i++ { // skip some basic frames.
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		// stacks += fmt.Sprintf("file: %s:%d func: %s\n", file, line, f.Name())
		builder.WriteString("file: " + file + ":" + strconv.Itoa(line) + " func: " + f.Name() + "\n")

	}
	return builder.String()
}

func newError(code ErrorCodeType, err error, msg ...string) error {
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}

	stacks := getStacks(skipStack)

	return &XDError{
		Code:     code,
		Message:  str,
		RTStacks: stacks,
		err:      err,
	}
}

func New(code ErrorCodeType, msg ...string) error {
	return newError(code, nil, msg...)
}

func ErrToXDError(err error, code ErrorCodeType) error {
	if err == nil {
		return nil
	}

	_, ok := err.(*XDError)
	if ok {
		return err
	}

	return newError(code, err, err.Error())
}

func WrapCode(err error, code ErrorCodeType) error {
	return ErrToXDError(err, code)
}

func Wrap(err error, msg ...string) error {
	if err == nil {
		return nil
	}

	_, ok := err.(*XDError)
	if ok {
		return err
	}

	var info string
	if len(msg) > 0 {
		// info = fmt.Sprintf("%s: %s", msg[0], err.Error())
		info = msg[0] + ": " + err.Error()
	} else {
		info = err.Error()
	}

	return newError(ErrorCodeUnknown, err, info)
}

func Code(err error) ErrorCodeType {
	xderror, ok := err.(*XDError)
	if !ok {
		return ErrorCodeUnknown
	}

	return xderror.Code
}

// newXDError 是对 newError 的兼容优化
func newXDError(code ErrorCodeType, err error, msgs ...string) *XDError {
	if err == nil {
		err = errors.New("")
	}
	xderr := &XDError{
		Code:     code,
		RTStacks: getStacks(4),
		err:      err,
	}
	if len(msgs) > 0 {
		xderr.Message = strings.Join(msgs, ",")
	} else {
		xderr.Message = err.Error()
	}
	return xderr
}

// NewXDErr 分配 IXDError 对象
func NewXDErr(code int, msg string) IXDError {
	return newXDError(ErrorCodeType(code), errors.New(msg), msg)
}

// NewXDError 分配 IXDError 对象
func NewXDError(code int, err error, msg ...string) IXDError {
	return newXDError(ErrorCodeType(code), err, msg...)
}

// NewBadRequestErr 分配带HTTP状态码(400)的 IXDError
func NewBadRequestErr(code int, err error, msg ...string) IXDError {
	return NewXDHTTPErr(http.StatusBadRequest, code, err, msg...)
}

// NewInternalErr 分配带HTTP状态码(500)的 IXDError
func NewInternalErr(code int, err error, msg ...string) IXDError {
	return NewXDHTTPErr(http.StatusInternalServerError, code, err, msg...)
}

// NewXDHTTPErr 分配带HTTP状态码的 IXDError
func NewXDHTTPErr(httpCode, code int, err error, msg ...string) IXDError {
	xderr := newXDError(ErrorCodeType(code), err, msg...)
	xderr.httpCode = httpCode
	return xderr
}
