package app

import (
	"fmt"
	"gin_kubelet/pkg/logger"
	"io"
	"log/slog"

	"gin_kubelet/config"
)

// App 应用容器，封装所有依赖
type App struct {
	Config *config.Config
	//DB         *gorm.DB
	//ClickHouse *sql.DB
	//Redis      *redis.Client
	Logger *slog.Logger

	// 内部资源，用于关闭
	logCloser io.Closer
}

// Option 应用选项函数
type Option func(*App) error

// New 创建新的应用实例
func New(cfg *config.Config, opts ...Option) (*App, error) {
	app := &App{
		Config: cfg,
	}

	for _, opt := range opts {
		if err := opt(app); err != nil {
			// 如果初始化失败，清理已创建的资源
			_ = app.Close()
			return nil, err
		}
	}

	return app, nil
}

// WithLogger 初始化日志
func WithLogger() Option {
	return func(a *App) error {
		log, closer, err := logger.New(&a.Config.Log)
		if err != nil {
			return fmt.Errorf("初始化日志失败: %w", err)
		}
		a.Logger = log
		a.logCloser = closer
		// 设置为全局默认日志（便于 slog 直接使用）
		slog.SetDefault(log)
		return nil
	}
}

//// WithDatabase 初始化数据库连接
//func WithDatabase() Option {
//	return func(a *App) error {
//		db, err := database.New(&a.Config.Database)
//		if err != nil {
//			return fmt.Errorf("初始化数据库失败: %w", err)
//		}
//		a.DB = db
//
//		// 执行分表迁移
//		if err := database.MigratePurgeURLResultTables(db); err != nil {
//			return fmt.Errorf("分表迁移失败: %w", err)
//		}
//
//		return nil
//	}
//}

//// WithRedis 初始化 Redis 连接
//func WithRedis() Option {
//	return func(a *App) error {
//		client, err := pkgredis.New(&a.Config.Redis)
//		if err != nil {
//			return fmt.Errorf("初始化 Redis 失败: %w", err)
//		}
//		a.Redis = client
//		return nil
//	}
//}

//// WithClickHouse 初始化 ClickHouse 连接
//func WithClickHouse() Option {
//	return func(a *App) error {
//		conn, err := clickhouse.New(&a.Config.ClickHouse)
//		if err != nil {
//			return fmt.Errorf("初始化 ClickHouse 失败: %w", err)
//		}
//		a.ClickHouse = conn
//		return nil
//	}
//}

// Close 关闭所有资源
func (a *App) Close() error {
	var errs []error

	//// 关闭数据库
	//if a.DB != nil {
	//	if err := database.Close(a.DB); err != nil {
	//		errs = append(errs, fmt.Errorf("关闭数据库失败: %w", err))
	//	}
	//}
	//
	//// 关闭 ClickHouse
	//if a.ClickHouse != nil {
	//	if err := clickhouse.Close(a.ClickHouse); err != nil {
	//		errs = append(errs, fmt.Errorf("关闭 ClickHouse 失败: %w", err))
	//	}
	//}
	//
	//// 关闭 Redis
	//if a.Redis != nil {
	//	if err := a.Redis.Close(); err != nil {
	//		errs = append(errs, fmt.Errorf("关闭 Redis 失败: %w", err))
	//	}
	//}

	// 关闭日志写入器
	if a.logCloser != nil {
		if err := a.logCloser.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭日志失败: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("关闭应用时发生错误: %v", errs)
	}

	return nil
}
