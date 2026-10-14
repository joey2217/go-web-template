// Package config 负责应用配置的结构定义与多环境加载。
package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultEnv = "dev"
	envPrefix  = "APP"
)

// Config 是应用的全量配置。
type Config struct {
	App       App       `mapstructure:"app"`
	Server    Server    `mapstructure:"server"`
	Database  Database  `mapstructure:"database"`
	Redis     Redis     `mapstructure:"redis"`
	RateLimit RateLimit `mapstructure:"ratelimit"`
	Cache     Cache     `mapstructure:"cache"`
	Metrics   Metrics   `mapstructure:"metrics"`
	Swagger   Swagger   `mapstructure:"swagger"`
	Log       Log       `mapstructure:"log"`
	CORS      CORS      `mapstructure:"cors"`
}

// App 应用元信息。
type App struct {
	Name  string `mapstructure:"name"`
	Env   string `mapstructure:"env"`
	Debug bool   `mapstructure:"debug"`
}

// Server HTTP 服务配置。
type Server struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// Database PostgreSQL 配置。
type Database struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	LogLevel        string        `mapstructure:"log_level"`
}

// Redis 配置。
type Redis struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// RateLimit 限流配置。
type RateLimit struct {
	Enabled   bool          `mapstructure:"enabled"`
	KeyPrefix string        `mapstructure:"key_prefix"`
	Rate      int           `mapstructure:"rate"`
	Burst     int           `mapstructure:"burst"`
	Window    time.Duration `mapstructure:"window"`
	FailOpen  bool          `mapstructure:"fail_open"`
}

// Cache 业务缓存配置（基于 Redis）。
type Cache struct {
	Enabled    bool          `mapstructure:"enabled"`
	KeyPrefix  string        `mapstructure:"key_prefix"`
	DefaultTTL time.Duration `mapstructure:"default_ttl"`
}

// Metrics Prometheus 指标配置。
type Metrics struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

// Swagger 文档开关配置。
type Swagger struct {
	Enabled bool `mapstructure:"enabled"`
}

// Log 日志配置。
type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// CORS 跨域配置。
type CORS struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
}

// IsProd 判断是否为生产环境。
func (c *Config) IsProd() bool { return c.App.Env == "prod" }

// Load 加载配置。
//
// 优先级：环境变量 > config.local.yaml > config.{env}.yaml > config.yaml > 代码默认值。
// 当显式传入 configFile 时，只读取该文件（仍支持环境变量覆盖）。
func Load(configFile string) (*Config, error) {
	v := viper.New()
	setDefaults(v)

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file %q: %w", configFile, err)
		}
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		if err := v.ReadInConfig(); err != nil {
			var notFound viper.ConfigFileNotFoundError
			if !errors.As(err, &notFound) {
				return nil, fmt.Errorf("read config: %w", err)
			}
		}

		env := v.GetString("app.env")
		if env == "" {
			env = defaultEnv
		}
		// 依次叠加环境配置与本地配置（后者优先级更高）。
		for _, name := range []string{fmt.Sprintf("config.%s", env), "config.local"} {
			if err := mergeFile(v, name); err != nil {
				return nil, err
			}
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if cfg.App.Env == "" {
		cfg.App.Env = defaultEnv
	}
	return cfg, nil
}

// mergeFile 在文件存在时将其配置合并进主配置，不存在则忽略。
func mergeFile(v *viper.Viper, name string) error {
	ov := viper.New()
	ov.SetConfigName(name)
	ov.SetConfigType("yaml")
	ov.AddConfigPath(".")
	ov.AddConfigPath("./config")
	if err := ov.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			return nil
		}
		return fmt.Errorf("read config %s: %w", name, err)
	}
	if err := v.MergeConfigMap(ov.AllSettings()); err != nil {
		return fmt.Errorf("merge config %s: %w", name, err)
	}
	return nil
}

// setDefaults 设置默认值，保证零配置也可启动。
func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "go-web-template")
	v.SetDefault("app.env", defaultEnv)
	v.SetDefault("app.debug", true)

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.shutdown_timeout", "10s")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "app")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 50)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "1h")
	v.SetDefault("database.log_level", "warn")

	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")

	v.SetDefault("ratelimit.enabled", true)
	v.SetDefault("ratelimit.key_prefix", "rl")
	v.SetDefault("ratelimit.rate", 20)
	v.SetDefault("ratelimit.burst", 40)
	v.SetDefault("ratelimit.window", "1s")
	v.SetDefault("ratelimit.fail_open", true)

	v.SetDefault("cache.enabled", true)
	v.SetDefault("cache.key_prefix", "gwt:")
	v.SetDefault("cache.default_ttl", "5m")

	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.path", "/metrics")

	v.SetDefault("swagger.enabled", true)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")

	v.SetDefault("cors.allow_origins", []string{"*"})
}
