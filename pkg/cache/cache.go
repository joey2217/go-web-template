// Package cache 提供缓存抽象与 Redis 实现。
//
// 模式：cache-aside（旁路缓存）。调用方负责维护一致性：
//   - 读：先查缓存，未命中回源并回填
//   - 写：更新/删除数据库后，主动 Delete 失效缓存项
//
// 不做 write-through / write-behind：复杂度高，对 CRUD 接口收益有限。
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// ErrCacheMiss 表示缓存未命中——这是正常分支，不是错误。
// 业务层用 errors.Is(err, cache.ErrCacheMiss) 判断并回源。
var ErrCacheMiss = errors.New("cache: miss")

// Cache 是缓存抽象。实现需保证：
//   - Get 在 key 不存在或已过期时返回 ErrCacheMiss
//   - Delete 幂等：key 不存在不报错
//   - DeleteByPattern 使用 SCAN + DEL，影响性能，避免热路径
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}

// Get 是带类型安全的读取辅助。未命中返回 (zero, false, nil)。
// 仅当底层发生真正的错误（网络 / 反序列化等）时返回 error。
func Get[T any](ctx context.Context, c Cache, key string) (T, bool, error) {
	var zero T
	raw, err := c.Get(ctx, key)
	if err != nil {
		if errors.Is(err, ErrCacheMiss) {
			return zero, false, nil
		}
		return zero, false, err
	}
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		return zero, false, err
	}
	return v, true, nil
}

// Set 是带 JSON 序列化的便捷写入。
func Set[T any](ctx context.Context, c Cache, key string, value T, ttl time.Duration) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, raw, ttl)
}

// Noop 是 cache.Cache 的空实现，用于未启用缓存或 Redis 不可用的场景。
// 所有 Get 返回 ErrCacheMiss，所有写操作直接成功（不持久化）。
type Noop struct{}

func (Noop) Get(context.Context, string) ([]byte, error) { return nil, ErrCacheMiss }
func (Noop) Set(context.Context, string, []byte, time.Duration) error {
	return nil
}
func (Noop) Delete(context.Context, ...string) error       { return nil }
func (Noop) DeleteByPattern(context.Context, string) error { return nil }
