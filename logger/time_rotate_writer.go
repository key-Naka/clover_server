package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// TimeRotateConfig 定义按时间切分的文件输出参数。
type TimeRotateConfig struct {
	Dir        string
	Prefix     string
	Interval   string
	MaxAgeDays int
	Enabled    bool
}

type timeRotateWriter struct {
	mu         sync.Mutex
	dir        string
	prefix     string
	interval   string
	maxAgeDays int
	enabled    bool
	currentKey string
	current    *os.File
}

func newTimeRotateWriter(cfg TimeRotateConfig) (*timeRotateWriter, error) {
	if cfg.Dir == "" {
		cfg.Dir = "logs"
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "app"
	}
	if cfg.Interval == "" {
		cfg.Interval = "daily"
	}
	if cfg.MaxAgeDays <= 0 {
		cfg.MaxAgeDays = 7
	}

	writer := &timeRotateWriter{
		dir:        cfg.Dir,
		prefix:     cfg.Prefix,
		interval:   normalizeInterval(cfg.Interval),
		maxAgeDays: cfg.MaxAgeDays,
		enabled:    cfg.Enabled,
	}

	if err := os.MkdirAll(writer.dir, 0o755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	if err := writer.rotateIfNeeded(time.Now()); err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *timeRotateWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	if err := w.rotateIfNeeded(now); err != nil {
		return 0, err
	}

	if w.current == nil {
		return 0, fmt.Errorf("日志文件未初始化")
	}

	return w.current.Write(p)
}

func (w *timeRotateWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.current == nil {
		return nil
	}

	err := w.current.Close()
	w.current = nil
	return err
}

func (w *timeRotateWriter) rotateIfNeeded(now time.Time) error {
	key := w.rotationKey(now)
	if w.current != nil && key == w.currentKey {
		return nil
	}

	if w.current != nil {
		if err := w.current.Close(); err != nil {
			return fmt.Errorf("关闭旧日志文件失败: %w", err)
		}
		w.current = nil
	}

	filename := filepath.Join(w.dir, w.fileName(now))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}

	w.current = file
	w.currentKey = key

	if err := w.cleanupExpired(now); err != nil {
		return err
	}

	return nil
}

func (w *timeRotateWriter) cleanupExpired(now time.Time) error {
	if w.maxAgeDays <= 0 {
		return nil
	}

	expireBefore := now.Add(-time.Duration(w.maxAgeDays) * 24 * time.Hour)
	pattern := filepath.Join(w.dir, fmt.Sprintf("%s-*.log", w.prefix))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("扫描历史日志失败: %w", err)
	}

	for _, filePath := range matches {
		info, statErr := os.Stat(filePath)
		if statErr != nil {
			return fmt.Errorf("读取日志文件状态失败: %w", statErr)
		}

		if info.ModTime().Before(expireBefore) {
			if removeErr := os.Remove(filePath); removeErr != nil && !os.IsNotExist(removeErr) {
				return fmt.Errorf("清理过期日志失败: %w", removeErr)
			}
		}
	}

	return nil
}

func (w *timeRotateWriter) rotationKey(now time.Time) string {
	if !w.enabled {
		return "static"
	}

	switch w.interval {
	case "hourly":
		return now.Format("2006010215")
	default:
		return now.Format("20060102")
	}
}

func (w *timeRotateWriter) fileName(now time.Time) string {
	if !w.enabled {
		return fmt.Sprintf("%s.log", w.prefix)
	}

	switch w.interval {
	case "hourly":
		return fmt.Sprintf("%s-%s.log", w.prefix, now.Format("2006-01-02-15"))
	default:
		return fmt.Sprintf("%s-%s.log", w.prefix, now.Format("2006-01-02"))
	}
}

func normalizeInterval(interval string) string {
	switch strings.ToLower(strings.TrimSpace(interval)) {
	case "hour", "hourly":
		return "hourly"
	case "day", "daily":
		return "daily"
	default:
		return "daily"
	}
}
