// Package response 提供统一的 HTTP 响应封装。
//
// 响应约定：
//   - 成功响应（2xx）直接返回业务数据，**不包信封**。
//     例：POST /users → 201 {"id":"...","username":"..."}
//   - 错误响应（4xx/5xx）使用 ErrorBody 信封，含业务错误码与提示。
//     例：409 → {"code":11001,"message":"邮箱已存在"}
//   - request_id 始终通过响应头 X-Request-Id 返回（由 RequestID 中间件写入），
//     成功与错误响应均不写入 body。
package response

import (
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/example/go-web-template/pkg/apperror"
)

// PageData 是分页响应的数据结构（直接返回，不包信封）。
type PageData struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// ErrorBody 是错误响应的统一信封。成功响应不使用本类型。
type ErrorBody struct {
	Code      int    `json:"code"`                 // 业务错误码
	Message   string `json:"message"`              // 提示信息
	RequestID string `json:"request_id,omitempty"` // 请求 ID，便于排查
}

// OK 返回 200 + 业务数据（无信封）。
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// Created 返回 201 + 业务数据（无信封）。
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

// Page 返回 200 + 分页数据（无信封）。
func Page(c *gin.Context, list any, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, PageData{List: list, Total: total, Page: page, PageSize: pageSize})
}

// NoContent 返回 204 无内容响应。
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Fail 将任意 error 渲染为统一错误响应（带信封）。
func Fail(c *gin.Context, err error) {
	appErr := apperror.From(err)
	c.AbortWithStatusJSON(appErr.HTTPStatus, ErrorBody{
		Code:      appErr.Code,
		Message:   appErr.Message,
		RequestID: requestid.Get(c),
	})
}
