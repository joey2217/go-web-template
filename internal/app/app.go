// Package app 负责依赖装配与应用生命周期管理。
package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/example/go-web-template/internal/config"
	"github.com/example/go-web-template/internal/handler"
	"github.com/example/go-web-template/internal/router"
	"github.com/example/go-web-template/pkg/cache"
	"github.com/example/go-web-template/pkg/database"
	"github.com/example/go-web-template/pkg/logger"
	appredis "github.com/example/go-web-template/pkg/redis"
)

// App 持有应用运行所需的全部长生命周期对象。
type App struct {
	cfg    *config.Config
	logger *zap.Logger
	db     *gorm.DB
	rdb    *goredis.Client
	server *http.Server
}

// New 初始化日志、数据库、Redis 并完成依赖装配。
func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log, err := logger.New(logger.Config{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
		Output: cfg.Log.Output,
	})
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	db, err := database.NewPostgres(database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Name:            cfg.Database.Name,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}, logger.NewGormLogger(log, gormLogLevel(cfg.Database.LogLevel)))
	if err != nil {
		log.Error("connect database failed", zap.Error(err))
		return nil, err
	}

	// Redis 在需要限流或缓存时连接，两者都关时也不连（避免无 Redis 环境下无法启动）。
	var rdb *goredis.Client
	if cfg.RateLimit.Enabled || cfg.Cache.Enabled {
		rdb, err = appredis.New(ctx, appredis.Config{
			Addr:         cfg.Redis.Addr,
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
			PoolSize:     cfg.Redis.PoolSize,
			DialTimeout:  cfg.Redis.DialTimeout,
			ReadTimeout:  cfg.Redis.ReadTimeout,
			WriteTimeout: cfg.Redis.WriteTimeout,
		})
		if err != nil {
			log.Error("connect redis failed", zap.Error(err))
			return nil, err
		}
	}

	// 构造缓存：未启用或 Redis 不可用时回退到 Noop，业务侧无感。
	var c cache.Cache = cache.Noop{}
	if cfg.Cache.Enabled && rdb != nil {
		c = cache.NewRedis(cache.RedisConfig{
			Client:    rdb,
			KeyPrefix: cfg.Cache.KeyPrefix,
		})
	}

	// 依赖装配：handler 聚合体（内部完成 repository → service → handler 串联）。
	handlers := handler.NewHandlers(db, c)

	engine := router.New(router.Deps{
		Config:   cfg,
		Logger:   log,
		DB:       db,
		Redis:    rdb,
		Handlers: handlers,
	})

	server := &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port)),
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &App{
		cfg:    cfg,
		logger: log,
		db:     db,
		rdb:    rdb,
		server: server,
	}, nil
}

// Run 启动 HTTP 服务并阻塞直到收到退出信号或发生致命错误。
func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		a.logger.Info("http server listening",
			zap.String("addr", a.server.Addr),
			zap.String("env", a.cfg.App.Env),
		)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		a.close()
		return err
	case <-ctx.Done():
		a.logger.Info("shutdown signal received, shutting down gracefully")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.Server.ShutdownTimeout)
		defer cancel()
		err := a.server.Shutdown(shutdownCtx)
		a.close()
		if err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
		a.logger.Info("server exited")
		return nil
	}
}

// close 释放数据库与 Redis 连接。
func (a *App) close() {
	if a.rdb != nil {
		if err := a.rdb.Close(); err != nil {
			a.logger.Warn("close redis failed", zap.Error(err))
		}
	}
	if a.db != nil {
		if sqlDB, err := a.db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				a.logger.Warn("close database failed", zap.Error(err))
			}
		}
	}
}

// gormLogLevel 将配置中的日志级别转换为 GORM 日志级别。
func gormLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "info":
		return gormlogger.Info
	default:
		return gormlogger.Warn
	}
}
