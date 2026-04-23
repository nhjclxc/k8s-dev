package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gin_kubelet/config"
	"gin_kubelet/pkg/app"
)

var (
	cfgFile   string
	globalApp *app.App
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "gin_kubelet",
	Short: "CDN 调度服务",
	Long:  `gin_kubelet 是部署在中心集群的 PCDN 调度服务，负责管理 PCDN 节点状态、计算调度权重，并将可用节点列表推送给边缘节点`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行命令失败: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// 持久化标志（所有子命令都可以使用）
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config/config.yaml", "配置文件路径")

	// 添加子命令
	rootCmd.AddCommand(httpCmd)
	rootCmd.AddCommand(cronCmd)
	rootCmd.AddCommand(startCmd)
}

// initApp 初始化应用
func initApp() (*app.App, error) {
	// 加载配置
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 创建应用实例，按顺序初始化各个组件
	application, err := app.New(cfg,
		app.WithLogger(), // 首先初始化日志
		//app.WithDatabase(),   // 初始化数据库
		//app.WithClickHouse(), // 初始化 ClickHouse
		//app.WithRedis(),      // 初始化 Redis
	)
	if err != nil {
		return nil, fmt.Errorf("初始化应用失败: %w", err)
	}

	globalApp = application
	return application, nil
}

// getApp 获取应用实例（用于子命令）
func getApp() *app.App {
	return globalApp
}
