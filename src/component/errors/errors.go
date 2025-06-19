package errors

import (
	"MetaFarmBankend/src/component/logger"
	"context"
	"runtime"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Stack   string `json:"stack"`
}

func New(code int, message string) *Error {
	_, file, line, _ := runtime.Caller(1)
	return &Error{
		Code:    code,
		Message: message,
		File:    file,
		Line:    line,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithStack() *Error {
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, false)
	e.Stack = string(buf[:n])
	return e
}

func (e *Error) LogError() {
	logger.Error("System error occurred",
		"code", e.Code,
		"message", e.Message,
		"file", e.File,
		"line", e.Line,
		"stack", e.Stack,
	)
}

func Recover(ctx context.Context) {
	if err := recover(); err != nil {
		switch v := err.(type) {
		case *Error:
			v.LogError()
		case error:
			_, file, line, _ := runtime.Caller(3)
			e := New(500, v.Error()).WithStack()
			e.File = file
			e.Line = line
			e.LogError()
		default:
			_, file, line, _ := runtime.Caller(3)
			e := New(500, "Unknown error").WithStack()
			e.File = file
			e.Line = line
			e.LogError()
		}
	}
}

func RecoverHandler(ctx context.Context, handler func(err *Error)) {
	if err := recover(); err != nil {
		switch v := err.(type) {
		case *Error:
			v.LogError()
			handler(v)
		case error:
			_, file, line, _ := runtime.Caller(3)
			e := New(500, v.Error()).WithStack()
			e.File = file
			e.Line = line
			e.LogError()
			handler(e)
		default:
			_, file, line, _ := runtime.Caller(3)
			e := New(500, "Unknown error").WithStack()
			e.File = file
			e.Line = line
			e.LogError()
			handler(e)
		}
	}
}

func Middleware(next func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		defer Recover(ctx)
		return next(ctx)
	}
}
