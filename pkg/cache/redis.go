package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// RedisConfig 是 RedisCache 的配置。
type RedisConfig struct {
	Client    *goredis.Client
	KeyPrefix string // 所有 key 自动加此前缀，便于多项目共享 Redis
}

// RedisCache 是基于 Redis 的 Cache 实现。
type RedisCache struct {
	client *goredis.Client
	prefix string
}

// NewRedis 创建 RedisCache。KeyPrefix 为空时默认 "cache:"。
func NewRedis(cfg RedisConfig) *RedisCache {
	prefix := cfg.KeyPrefix
	if prefix == "" {
		prefix = "cache:"
	}
	return &RedisCache{client: cfg.Client, prefix: prefix}
}

func (c *RedisCache) key(k string) string { return c.prefix + k }

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	b, err := c.client.Get(ctx, c.key(key)).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("cache get %q: %w", key, err)
	}
	return b, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := c.client.Set(ctx, c.key(key), value, ttl).Err(); err != nil {
		return fmt.Errorf("cache set %q: %w", key, err)
	}
	return nil
}

func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	full := make([]string, len(keys))
	for i, k := range keys {
		full[i] = c.key(k)
	}
	if err := c.client.Del(ctx, full...).Err(); err != nil {
		return fmt.Errorf("cache delete: %w", err)
	}
	return nil
}

// DeleteByPattern 使用 SCAN 分批删除匹配 key，避免 KEYS 阻塞 Redis。
// 单批上限 200，影响性能，仅用于写后批量失效。
func (c *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	fullPattern := c.prefix + pattern
	iter := c.client.Scan(ctx, 0, fullPattern, 200).Iterator()
	batch := make([]string, 0, 200)
	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		if err := c.client.Del(ctx, batch...).Err(); err != nil {
			return fmt.Errorf("cache delete pattern: %w", err)
		}
		batch = batch[:0]
		return nil
	}
	for iter.Next(ctx) {
		batch = append(batch, iter.Val())
		if len(batch) >= 200 {
			if err := flush(); err != nil {
				return err
			}
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("cache scan: %w", err)
	}
	return flush()
}
