// Package apperror 定义统一的业务错误类型与错误码。
//
// 错误码分段：
//
//	0      成功
//	10xxx  通用 / 参数
//	11xxx  认证鉴权
//	12xxx  业务
//	13xxx  限流
//	50xxx  系统
package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// 业务错误码。
const (
	CodeOK = 0

	CodeBadRequest   = 10001
	CodeUnauthorized = 11001
	CodeForbidden    = 11002
	CodeNotFound     = 12001
	CodeConflict     = 12002
	CodeTooManyReq   = 13001
	CodeInternal     = 50000
)

// AppError 表示一个可安全暴露给客户端的业务错误。
type AppError struct {
	Code       int    // 业务错误码
	HTTPStatus int    // HTTP 状态码
	Message    string // 面向客户端的提示信息
	Cause      error  // 内部原因，仅用于日志，不对外暴露
}

// Error 实现 error 接口。
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持 errors.Is / errors.As 链式判断。
func (e *AppError) Unwrap() error { return e.Cause }

// WithCause 返回携带内部原因的错误副本。
func (e *AppError) WithCause(err error) *AppError {
	clone := *e
	clone.Cause = err
	return &clone
}

// WithMessage 返回替换提示信息的错误副本。
func (e *AppError) WithMessage(msg string) *AppError {
	clone := *e
	clone.Message = msg
	return &clone
}

// New 构造一个业务错误。
func New(code, httpStatus int, message string) *AppError {
	return &AppError{Code: code, HTTPStatus: httpStatus, Message: message}
}

// BadRequest 参数错误（400）。
func BadRequest(message string) *AppError {
	return New(CodeBadRequest, http.StatusBadRequest, message)
}

// Unauthorized 未认证（401）。
func Unauthorized(message string) *AppError {
	return New(CodeUnauthorized, http.StatusUnauthorized, message)
}

// Forbidden 无权限（403）。
func Forbidden(message string) *AppError {
	return New(CodeForbidden, http.StatusForbidden, message)
}

// NotFound 资源不存在（404）。
func NotFound(message string) *AppError {
	return New(CodeNotFound, http.StatusNotFound, message)
}

// Conflict 资源冲突（409）。
func Conflict(message string) *AppError {
	return New(CodeConflict, http.StatusConflict, message)
}

// TooManyRequests 请求过于频繁（429）。
func TooManyRequests(message string) *AppError {
	return New(CodeTooManyReq, http.StatusTooManyRequests, message)
}

// Internal 内部错误（500）。
func Internal(message string) *AppError {
	return New(CodeInternal, http.StatusInternalServerError, message)
}

// From 将任意 error 归一化为 *AppError；未知错误统一视为内部错误且不暴露细节。
func From(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return Internal("服务器内部错误").WithCause(err)
}
