package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/example/go-web-template/pkg/apperror"
	"github.com/example/go-web-template/pkg/logger"
	"github.com/example/go-web-template/pkg/response"
)

// ErrorHandler 在链路末端统一渲染错误。
//
// handler 通过 c.Error(err) 上报错误，本中间件将其映射为统一的响应体；
// 5xx 会额外记录错误日志（自动携带 request_id），便于排查。
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		appErr := apperror.From(c.Errors.Last().Err)
		if appErr.HTTPStatus >= 500 {
			log := logger.FromGin(c)
			log.Error("request failed",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Error(appErr),
			)
		}
		response.Fail(c, appErr)
	}
}
