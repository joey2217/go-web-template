package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type loggerCtxKey struct{}

// WithContext 返回携带指定 logger 的派生 context。
// 通常在请求入口（中间件）调用一次，下游通过 FromContext 取出。
func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, l)
}

// FromContext 从 context 取出 logger；若未设置则返回 no-op logger，
// 以保证调用方不会出现空指针，可以无条件调用 log.Info / log.Error 等方法。
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger); ok && l != nil {
		return l
	}
	return zap.NewNop()
}

// FromGin 从 Gin 上下文取出请求级 logger。
// 这是 handler 中最常用的入口：log := logger.FromGin(c)。
// 取出的 logger 已经携带 request_id 等请求级字段，无需再手动追加。
func FromGin(c *gin.Context) *zap.Logger {
	return FromContext(c.Request.Context())
}
