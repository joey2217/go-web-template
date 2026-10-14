// Package middleware 提供 Gin 中间件。
package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/example/go-web-template/pkg/logger"
)

// RequestID 为每个请求生成唯一 ID 写入响应头 X-Request-Id，
// 并基于 baseLog 创建带 request_id 字段的子 logger 放入请求 context。
//
// 下游中间件与 handler 通过 logger.FromGin(c) 即可获得自动携带 request_id 的 logger，
// 无需在每条日志上手动 zap.String("request_id", ...)。
func RequestID(baseLog *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestid.New()(c)
		id := requestid.Get(c)
		child := baseLog.With(zap.String("request_id", id))
		c.Request = c.Request.WithContext(logger.WithContext(c.Request.Context(), child))
	}
}
