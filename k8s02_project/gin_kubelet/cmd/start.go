package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"gin_kubelet/internal/server"
)

// startCmd 启动所有服务命令
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动所有服务",
	Long:  `启动所有服务，包括 HTTP API、gRPC 调度服务、队列任务和检测定时任务`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 初始化应用
		application, err := initApp()
		if err != nil {
			return err
		}

		slog.Info("启动所有服务...")

		// 创建服务器实例
		httpServer := server.NewHTTPServer(application)
		cronServer := server.NewCronServer(application)

		// 使用 WaitGroup 等待所有服务启动
		var wg sync.WaitGroup

		// 启动 HTTP 服务
		wg.Go(func() {
			if err := httpServer.Start(); err != nil {
				slog.Error("HTTP 服务启动失败", "error", err)
			}
		})

		// 启动定时任务服务
		wg.Go(func() {
			if err := cronServer.Start(); err != nil {
				slog.Error("定时任务服务启动失败", "error", err)
			}
		})

		// 监听退出信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		slog.Info("正在优雅关闭服务...")

		// 创建关闭上下文，设置超时时间
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			slog.Error("关闭 HTTP 服务失败", "error", err)
		}

		if err := cronServer.Stop(ctx); err != nil {
			slog.Error("关闭定时任务服务失败", "error", err)
		}

		slog.Info("所有服务已关闭")

		// 关闭应用（包括数据库、Redis、Telemetry、日志等）
		if err := application.Close(); err != nil {
			slog.Error("关闭应用资源失败", "error", err)
		}

		return nil
	},
}
