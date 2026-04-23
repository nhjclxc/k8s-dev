package cmd

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	"gin_kubelet/internal/server"
)

var (
	taskName string
	taskArgs []string
)

// cronCmd 定时任务服务命令
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "启动定时任务服务",
	Long:  `启动定时任务服务，执行心跳超时检测等定时任务`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 初始化应用
		application, err := initApp()
		if err != nil {
			return err
		}
		defer application.Close()

		// 如果指定了 task 参数，则立即执行该任务
		if taskName != "" {
			slog.Info("立即执行任务", "task", taskName, "args", taskArgs)
			cronServer := server.NewCronServer(application)
			return cronServer.RunTask(context.Background(), taskName, taskArgs)
		}

		slog.Info("启动定时任务服务...")

		// 创建并启动定时任务服务器
		cronServer := server.NewCronServer(application)
		return cronServer.Start()
	},
}

func init() {
	cronCmd.Flags().StringVar(&taskName, "task", "", "立即执行指定的任务（用于测试）")
	cronCmd.Flags().StringArrayVar(&taskArgs, "args", []string{}, "任务参数")
}
