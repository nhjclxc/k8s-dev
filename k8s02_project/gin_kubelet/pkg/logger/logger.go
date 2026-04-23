package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"

	"gin_kubelet/config"
)

// New 创建新的日志实例
func New(cfg *config.LogConfig) (*slog.Logger, io.Closer, error) {
	// 解析日志级别
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var writer io.Writer
	var closer io.Closer

	// 根据配置决定输出位置
	if cfg.Output == "file" {
		// 确保日志目录存在
		dir := filepath.Dir(cfg.FilePath)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, nil, err
			}
		}

		// 使用 lumberjack 进行日志切割
		lj := &lumberjack.Logger{
			Filename:   cfg.FilePath,   // 日志文件路径
			MaxSize:    cfg.MaxSize,    // 单个文件最大大小（MB）
			MaxBackups: cfg.MaxBackups, // 保留的旧文件数量
			MaxAge:     cfg.MaxAge,     // 保留的最大天数
			Compress:   cfg.Compress,   // 是否压缩旧文件
			LocalTime:  true,           // 使用本地时间
		}
		writer = lj
		closer = lj
	} else {
		// 输出到标准输出
		writer = os.Stdout
		closer = nil
	}

	// 根据格式创建 Handler
	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	logger := slog.New(handler)

	return logger, closer, nil
}
