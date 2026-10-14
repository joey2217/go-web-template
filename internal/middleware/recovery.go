package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/example/go-web-template/pkg/apperror"
	"github.com/example/go-web-template/pkg/logger"
	"github.com/example/go-web-template/pkg/response"
)

// Recovery 捕获 panic，记录堆栈并返回统一的 500 响应。
// 日志由 RequestID 注入的请求级 logger 输出，自动携带 request_id。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log := logger.FromGin(c)
				log.Error("panic recovered",
					zap.Any("panic", rec),
					zap.String("path", c.Request.URL.Path),
					zap.Stack("stack"),
				)
				response.Fail(c, apperror.Internal("服务器内部错误"))
			}
		}()
		c.Next()
	}
}
