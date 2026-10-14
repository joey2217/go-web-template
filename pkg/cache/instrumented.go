package cache

import (
	"context"
	"errors"
	"time"

	"github.com/example/go-web-template/pkg/metrics"
)

// InstrumentedCache 装饰任意 Cache，自动上报操作计数到 Prometheus。
// 推荐用法：cache.NewInstrumented(cache.NewRedis(...))
// 让 RedisCache / Noop 保持纯净，不知道 metrics 存在（单一职责）。
type InstrumentedCache struct {
	inner Cache
}

// NewInstrumented 包装任意 Cache，自动计数所有操作。
func NewInstrumented(inner Cache) *InstrumentedCache {
	return &InstrumentedCache{inner: inner}
}

// Get: hit | miss | error
//   - hit:   缓存命中
//   - miss:  ErrCacheMiss（业务正常分支）
//   - error: 其它错误（网络 / 反序列化等）
func (c *InstrumentedCache) Get(ctx context.Context, key string) ([]byte, error) {
	b, err := c.inner.Get(ctx, key)
	var result string
	switch {
	case err == nil:
		result = "hit"
	case errors.Is(err, ErrCacheMiss):
		result = "miss"
	default:
		result = "error"
	}
	metrics.CacheOps.WithLabelValues("get", result).Inc()
	return b, err
}

// Set / Delete / DeleteByPattern: ok | error
func (c *InstrumentedCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := c.inner.Set(ctx, key, value, ttl)
	recordOp("set", err == nil)
	return err
}

func (c *InstrumentedCache) Delete(ctx context.Context, keys ...string) error {
	err := c.inner.Delete(ctx, keys...)
	recordOp("delete", err == nil)
	return err
}

func (c *InstrumentedCache) DeleteByPattern(ctx context.Context, pattern string) error {
	err := c.inner.DeleteByPattern(ctx, pattern)
	recordOp("delete_pattern", err == nil)
	return err
}

func recordOp(op string, ok bool) {
	result := "ok"
	if !ok {
		result = "error"
	}
	metrics.CacheOps.WithLabelValues(op, result).Inc()
}
