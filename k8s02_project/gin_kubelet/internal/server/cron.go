package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/robfig/cron/v3"

	internalCron "gin_kubelet/internal/cron"
	"gin_kubelet/pkg/app"
)

// CronServer 定时任务服务器
type CronServer struct {
	cron *cron.Cron
	app  *app.App
}

// NewCronServer 创建定时任务服务器实例
func NewCronServer(application *app.App) *CronServer {
	// 创建 cron 实例
	c := cron.New(
		cron.WithSeconds(), // 支持秒级定时
		cron.WithChain( // 添加中间件
			cron.Recover(cron.DefaultLogger), // 恢复 panic
		),
	)

	return &CronServer{
		cron: c,
		app:  application,
	}
}

// Start 启动定时任务服务器
func (s *CronServer) Start() error {
	cfg := s.app.Config

	// 注册所有任务
	for _, taskCfg := range cfg.Cron.Tasks {
		if !taskCfg.Enabled {
			slog.Info("跳过未启用的任务", "task", taskCfg.Name)
			continue
		}

		task := internalCron.GetTask(taskCfg.Name, s.app)
		if task == nil {
			slog.Warn("未找到任务", "task", taskCfg.Name)
			continue
		}

		// 捕获变量，避免闭包问题
		taskInstance := task
		taskName := taskCfg.Name

		// 添加任务
		_, err := s.cron.AddFunc(taskCfg.Spec, func() {
			s.executeTask(taskInstance, taskName)
		})

		if err != nil {
			return fmt.Errorf("添加定时任务失败: %s, %w", taskCfg.Name, err)
		}

		slog.Info("注册定时任务", "task", taskCfg.Name, "spec", taskCfg.Spec)
	}

	slog.Info("定时任务服务器启动")
	s.cron.Start()

	// 阻塞等待
	select {}
}

// executeTask 执行任务，包含 tracing 支持
func (s *CronServer) executeTask(task internalCron.Task, taskName string) {
	// 创建 root span（如果启用了 tracing）
	var ctx context.Context = context.Background()

	slog.Info("开始执行定时任务", "task", taskName)

	if err := task.Execute(ctx); err != nil {
		slog.Error("定时任务执行失败", "task", taskName, "error", err)
	} else {
		slog.Info("定时任务执行成功", "task", taskName)
	}
}

// Stop 停止定时任务服务器
func (s *CronServer) Stop(ctx context.Context) error {
	slog.Info("正在关闭定时任务服务器...")
	if s.cron != nil {
		stopCtx := s.cron.Stop()
		<-stopCtx.Done()
	}
	return nil
}

// RunTask 立即执行指定任务（用于测试）
func (s *CronServer) RunTask(ctx context.Context, taskName string, args []string) error {
	task := internalCron.GetTask(taskName, s.app)
	if task == nil {
		return fmt.Errorf("未找到任务: %s", taskName)
	}

	slog.Info("立即执行任务", "task", taskName, "args", args)
	if err := task.Execute(ctx); err != nil {
		slog.Error("任务执行失败", "task", taskName, "error", err)
		return err
	}

	slog.Info("任务执行成功", "task", taskName)
	return nil
}
