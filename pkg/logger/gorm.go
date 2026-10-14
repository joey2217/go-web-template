package logger

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// GormLogger 是接入 Zap 的 GORM 日志适配器。
type GormLogger struct {
	Zap           *zap.Logger
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
}

// NewGormLogger 创建 GORM 日志适配器。
func NewGormLogger(z *zap.Logger, level gormlogger.LogLevel) *GormLogger {
	return &GormLogger{
		Zap:           z,
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      level,
	}
}

// LogMode 实现 gormlogger.Interface。
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	clone := *l
	clone.LogLevel = level
	return &clone
}

// Info 实现 gormlogger.Interface。
func (l *GormLogger) Info(_ context.Context, msg string, data ...any) {
	if l.LogLevel >= gormlogger.Info {
		l.Zap.Sugar().Infof(msg, data...)
	}
}

// Warn 实现 gormlogger.Interface。
func (l *GormLogger) Warn(_ context.Context, msg string, data ...any) {
	if l.LogLevel >= gormlogger.Warn {
		l.Zap.Sugar().Warnf(msg, data...)
	}
}

// Error 实现 gormlogger.Interface。
func (l *GormLogger) Error(_ context.Context, msg string, data ...any) {
	if l.LogLevel >= gormlogger.Error {
		l.Zap.Sugar().Errorf(msg, data...)
	}
}

// Trace 实现 gormlogger.Interface。
func (l *GormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.String("source", utils.FileWithLineNum()),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}

	switch {
	case err != nil && !errors.Is(err, gormlogger.ErrRecordNotFound) && l.LogLevel >= gormlogger.Error:
		l.Zap.Error("gorm query error", append(fields, zap.Error(err))...)
	case l.SlowThreshold > 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		l.Zap.Warn("gorm slow query", fields...)
	case l.LogLevel >= gormlogger.Info:
		l.Zap.Debug("gorm query", fields...)
	}
}
