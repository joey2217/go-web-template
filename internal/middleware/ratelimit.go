package middleware

import (
	"strconv"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/example/go-web-template/internal/config"
	"github.com/example/go-web-template/pkg/apperror"
	"github.com/example/go-web-template/pkg/response"
)

// RateLimit 基于 Redis 的分布式限流（GCRA）。
//
// 限流维度为客户端 IP；接入鉴权后可按用户身份细分。
// 当 Redis 不可用时，依据 fail_open 决定放行还是拒绝。
func RateLimit(cfg config.RateLimit, rdb *goredis.Client, log *zap.Logger) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	limiter := redis_rate.NewLimiter(rdb)
	limit := redis_rate.Limit{
		Rate:   cfg.Rate,
		Burst:  cfg.Burst,
		Period: cfg.Window,
	}
	prefix := cfg.KeyPrefix
	if prefix == "" {
		prefix = "rl"
	}

	return func(c *gin.Context) {
		key := prefix + ":" + c.ClientIP()
		res, err := limiter.Allow(c.Request.Context(), key, limit)
		if err != nil {
			if cfg.FailOpen {
				log.Warn("rate limiter unavailable, fail open",
					zap.Error(err),
					zap.String("request_id", requestid.Get(c)),
				)
				c.Next()
				return
			}
			response.Fail(c, apperror.Internal("限流服务不可用").WithCause(err))
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(res.Limit.Rate))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))

		if res.Allowed == 0 {
			retryAfter := res.RetryAfter
			if retryAfter <= 0 {
				retryAfter = time.Second
			}
			c.Header("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
			response.Fail(c, apperror.TooManyRequests("请求过于频繁，请稍后再试"))
			return
		}

		c.Next()
	}
}
