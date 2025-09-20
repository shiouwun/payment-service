package errors

import (
	"fmt"
	"runtime"
)

type AppError struct {
	Message string
	Code    string
	Cause   error
	File    string
	Line    int
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v (at %s:%d)", e.Message, e.Cause, e.File, e.Line)
	}
	return fmt.Sprintf("%s (at %s:%d)", e.Message, e.File, e.Line)
}

func New(message string) *AppError {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Message: message,
		File:    file,
		Line:    line,
	}
}

func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Message: message,
		Cause:   err,
		File:    file,
		Line:    line,
	}
}

func WithCode(err error, code string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		appErr.Code = code
		return appErr
	}
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Message: err.Error(),
		Code:    code,
		Cause:   err,
		File:    file,
		Line:    line,
	}
}