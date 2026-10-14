// Package logger 封装 Zap 日志，并提供 GORM 日志适配。
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 是日志初始化配置。
type Config struct {
	Level  string // debug | info | warn | error
	Format string // json | console
	Output string // stdout 或文件路径
}

// New 按配置构建 *zap.Logger。
func New(cfg Config) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if cfg.Level != "" {
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			level = zapcore.InfoLevel
		}
	}

	var zc zap.Config
	if cfg.Format == "console" {
		zc = zap.NewDevelopmentConfig()
		zc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zc = zap.NewProductionConfig()
	}
	zc.Level = zap.NewAtomicLevelAt(level)

	if cfg.Output != "" && cfg.Output != "stdout" {
		zc.OutputPaths = []string{cfg.Output}
		zc.ErrorOutputPaths = []string{cfg.Output}
	}

	return zc.Build()
}
