package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Config 定义日志初始化参数。
type Config struct {
	Env     string
	Level   string
	Format  string
	Dir     string
	Prefix  string
	Console bool
	File    bool
	Rotate  RotateConfig
}

// RotateConfig 定义日志滚动策略。
type RotateConfig struct {
	Enabled    bool
	Interval   string
	MaxAgeDays int
}

// Init 初始化全局默认日志器。
func Init(cfg Config) error {
	writer, closeFn, err := buildWriter(cfg)
	if err != nil {
		return fmt.Errorf("构建日志输出失败: %w", err)
	}

	if closeFn != nil {
		registerClose(closeFn)
	}

	handlerOptions := &slog.HandlerOptions{
		Level:     parseLevel(cfg.Level, cfg.Env),
		AddSource: strings.EqualFold(cfg.Env, "dev"),
	}

	logger := slog.New(newHandler(writer, cfg.Format, cfg.Env, handlerOptions))
	slog.SetDefault(logger)

	return nil
}

func newHandler(writer io.Writer, format string, env string, opts *slog.HandlerOptions) slog.Handler {
	if strings.EqualFold(format, "json") && !strings.EqualFold(env, "dev") {
		return slog.NewJSONHandler(writer, opts)
	}

	if strings.EqualFold(format, "json") {
		return slog.NewJSONHandler(writer, opts)
	}

	return slog.NewTextHandler(writer, opts)
}

func buildWriter(cfg Config) (io.Writer, func() error, error) {
	writers := make([]io.Writer, 0, 2)
	var closeFn func() error

	if cfg.Console {
		writers = append(writers, os.Stdout)
	}

	if cfg.File {
		fileWriter, err := newTimeRotateWriter(TimeRotateConfig{
			Dir:        cfg.Dir,
			Prefix:     cfg.Prefix,
			Interval:   cfg.Rotate.Interval,
			MaxAgeDays: cfg.Rotate.MaxAgeDays,
			Enabled:    cfg.Rotate.Enabled,
		})
		if err != nil {
			return nil, nil, err
		}

		writers = append(writers, fileWriter)
		closeFn = fileWriter.Close
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	if len(writers) == 1 {
		return writers[0], closeFn, nil
	}

	return io.MultiWriter(writers...), closeFn, nil
}

func parseLevel(level string, env string) slog.Leveler {
	if strings.TrimSpace(level) == "" {
		if strings.EqualFold(env, "dev") {
			return slog.LevelDebug
		}
		return slog.LevelInfo
	}

	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		if strings.EqualFold(env, "dev") {
			return slog.LevelDebug
		}
		return slog.LevelInfo
	}
}

var closeCallbacks []func() error

func registerClose(closeFn func() error) {
	closeCallbacks = append(closeCallbacks, closeFn)
}

// Close 关闭日志底层资源。
func Close() error {
	var firstErr error
	for _, closeFn := range closeCallbacks {
		if err := closeFn(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	closeCallbacks = nil
	return firstErr
}
