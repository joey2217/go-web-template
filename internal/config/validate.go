package config

import (
	"errors"
	"fmt"
	"strings"
)

// Validate 校验关键配置项的合法性。
func (c *Config) Validate() error {
	var errs []error

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		errs = append(errs, fmt.Errorf("server.port 必须在 1-65535 之间，当前为 %d", c.Server.Port))
	}
	if c.Database.Host == "" {
		errs = append(errs, errors.New("database.host 不能为空"))
	}
	if c.Database.Name == "" {
		errs = append(errs, errors.New("database.name 不能为空"))
	}
	if c.Metrics.Enabled {
		if c.Metrics.Username == "" {
			errs = append(errs, errors.New("启用 metrics 时 metrics.username 不能为空"))
		}
		if strings.TrimSpace(c.Metrics.Password) == "" {
			errs = append(errs, errors.New("启用 metrics 时 metrics.password 不能为空"))
		}
	}
	if c.RateLimit.Enabled {
		if c.Redis.Addr == "" {
			errs = append(errs, errors.New("启用限流时 redis.addr 不能为空"))
		}
		if c.RateLimit.Rate <= 0 {
			errs = append(errs, errors.New("ratelimit.rate 必须大于 0"))
		}
		if c.RateLimit.Window <= 0 {
			errs = append(errs, errors.New("ratelimit.window 必须大于 0"))
		}
	}

	return errors.Join(errs...)
}
