package cron

import (
	"context"

	"gin_kubelet/pkg/app"
)

// Task 定时任务接口
type Task interface {
	// Name 返回任务名称
	Name() string
	// Execute 执行任务
	Execute(ctx context.Context) error
}

// TaskFactory 任务工厂函数类型
type TaskFactory func(app *app.App) Task

// TaskRegistry 任务注册表
var TaskRegistry = map[string]TaskFactory{
	"example_task": NewExampleTask,
}

// GetTask 根据名称获取任务
func GetTask(name string, application *app.App) Task {
	if factory, ok := TaskRegistry[name]; ok {
		return factory(application)
	}
	return nil
}
