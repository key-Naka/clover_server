package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
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
func Init(cfg Config) {
	writer, closeFn, err := buildWriter(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		os.Exit(1)
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
}

func newHandler(writer io.Writer, format string, env string, opts *slog.HandlerOptions) slog.Handler {
	if strings.EqualFold(format, "json") {
		return slog.NewJSONHandler(writer, opts)
	}

	return &customTextHandler{
		opts:   opts,
		writer: writer,
	}
}

type customTextHandler struct {
	opts   *slog.HandlerOptions
	writer io.Writer
	attrs  []slog.Attr
	groups []string
}

func (h *customTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts != nil && h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *customTextHandler) Handle(ctx context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006-01-02 15:04:05")
	levelStr := r.Level.String()

	sourceStr := ""
	if h.opts != nil && h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			sourceStr = fmt.Sprintf("%s:%d ", f.File, f.Line)
		}
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("[%s] [%s] %s%s", timeStr, levelStr, sourceStr, r.Message))

	for _, attr := range h.attrs {
		builder.WriteString(fmt.Sprintf(" %s=%v", attr.Key, attr.Value.Any()))
	}

	r.Attrs(func(a slog.Attr) bool {
		builder.WriteString(fmt.Sprintf(" %s=%v", a.Key, a.Value.Any()))
		return true
	})

	builder.WriteString("\n")
	_, err := h.writer.Write([]byte(builder.String()))
	return err
}

func (h *customTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := append([]slog.Attr{}, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &customTextHandler{
		opts:   h.opts,
		writer: h.writer,
		attrs:  newAttrs,
		groups: h.groups,
	}
}

func (h *customTextHandler) WithGroup(name string) slog.Handler {
	newGroups := append([]string{}, h.groups...)
	newGroups = append(newGroups, name)
	return &customTextHandler{
		opts:   h.opts,
		writer: h.writer,
		attrs:  h.attrs,
		groups: newGroups,
	}
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
