package core

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type SlogGormLogger struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  logger.LogLevel
}

// NewSlogGormLogger 创建一个基于 slog 的 GORM logger
func NewSlogGormLogger() logger.Interface {
	return &SlogGormLogger{
		SlowThreshold:             time.Second, // 将慢 SQL 阈值设为 1 秒
		IgnoreRecordNotFoundError: true,        // 忽略 ErrRecordNotFound 错误
		LogLevel:                  logger.Warn, // 默认日志级别
	}
}

func (l *SlogGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *SlogGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		slog.Info(msg, slog.Any("data", data))
	}
}

func (l *SlogGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		slog.Warn(msg, slog.Any("data", data))
	}
}

func (l *SlogGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		slog.Error(msg, slog.Any("data", data))
	}
}

func (l *SlogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		slog.Error("GORM ERROR",
			slog.String("sql", sql),
			slog.Any("err", err),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
			slog.String("src", utils.FileWithLineNum()),
		)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slog.Warn("GORM SLOW SQL",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
			slog.String("src", utils.FileWithLineNum()),
		)
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		slog.Debug("GORM SQL",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
			slog.String("src", utils.FileWithLineNum()),
		)
	}
}
