package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

// 配置结构体
type Config struct {
	App  AppConfig  `mapstructure:"app"`
	HTTP HTTPConfig `mapstructure:"http"`
	Log  LogConfig  `mapstructure:"log"`
	Cron CronConfig `mapstructure:"cron"`
}

// HTTPConfig HTTP 服务配置
type HTTPConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// CronConfig 定时任务配置
type CronConfig struct {
	Timezone string     `mapstructure:"timezone"`
	Tasks    []CronTask `mapstructure:"tasks"`
}

// CronTask 单个定时任务配置
type CronTask struct {
	Name    string `mapstructure:"name"`
	Spec    string `mapstructure:"spec"`
	Enabled bool   `mapstructure:"enabled"`
}

// AppConfig 应用程序配置
type AppConfig struct {
	Name    string `mapstructure:"name" json:"name"`
	Env     string `mapstructure:"env" json:"env"`
	Version string `mapstructure:"version" json:"version"`
	Debug   bool   `mapstructure:"debug"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`       // 日志级别: debug, info, warn, error
	Format     string `mapstructure:"format"`      // 日志格式: json, text
	Output     string `mapstructure:"output"`      // 输出方式: stdout, file
	FilePath   string `mapstructure:"file_path"`   // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 单个文件最大大小（MB）
	MaxBackups int    `mapstructure:"max_backups"` // 保留的旧文件数量
	MaxAge     int    `mapstructure:"max_age"`     // 保留的最大天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩旧文件
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 启用环境变量覆盖，前缀为 gin_kubelet，例如 GIN_KUBELET
	v.SetEnvPrefix("GIN_KUBELET")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}
