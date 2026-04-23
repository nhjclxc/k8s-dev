package cron

import (
	"context"
	"gin_kubelet/pkg/app"
	"log/slog"
)

// ExampleTask 示例任务
type ExampleTask struct {
	app *app.App
}

func NewExampleTask(application *app.App) Task {
	return &ExampleTask{
		app: application,
	}
}

func (t *ExampleTask) Name() string {
	return "example_task"
}

func (t *ExampleTask) Execute(ctx context.Context) error {
	slog.Info("执行示例任务", "task", t.Name())
	// 这里添加具体的任务逻辑
	// 可以使用 t.app.DB 访问数据库
	// 可以使用 t.app.Redis 访问 Redis
	return nil
}
