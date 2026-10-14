// Package response 提供统一的 HTTP 响应封装。
package response

import (
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/example/go-web-template/pkg/apperror"
)

// Body 是统一响应信封。
type Body struct {
	Code      int    `json:"code"`                 // 业务错误码，0 表示成功
	Message   string `json:"message"`              // 提示信息
	Data      any    `json:"data,omitempty"`       // 业务数据
	RequestID string `json:"request_id,omitempty"` // 请求 ID，便于排查
}

// PageData 是分页响应体。
type PageData struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// OK 返回 200 成功响应。
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:      apperror.CodeOK,
		Message:   "ok",
		Data:      data,
		RequestID: requestid.Get(c),
	})
}

// Created 返回 201 成功响应。
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Body{
		Code:      apperror.CodeOK,
		Message:   "ok",
		Data:      data,
		RequestID: requestid.Get(c),
	})
}

// Page 返回 200 分页响应。
func Page(c *gin.Context, list any, total int64, page, pageSize int) {
	OK(c, PageData{List: list, Total: total, Page: page, PageSize: pageSize})
}

// NoContent 返回 204 无内容响应。
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Fail 将任意 error 渲染为统一错误响应。
func Fail(c *gin.Context, err error) {
	appErr := apperror.From(err)
	c.AbortWithStatusJSON(appErr.HTTPStatus, Body{
		Code:      appErr.Code,
		Message:   appErr.Message,
		RequestID: requestid.Get(c),
	})
}
