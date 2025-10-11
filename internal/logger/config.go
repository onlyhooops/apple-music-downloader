package logger

import (
	"fmt"
	"io"
	"os"
)

// Config 日志配置
type Config struct {
	Level         string `yaml:"level"`          // 日志等级: debug/info/warn/error
	Output        string `yaml:"output"`         // 输出目标: stdout/stderr/文件路径
	ShowTimestamp bool   `yaml:"show_timestamp"` // 是否显示时间戳
}

// InitFromConfig 从配置初始化全局logger
func InitFromConfig(cfg Config) error {
	// 解析日志等级
	level := ParseLevel(cfg.Level)
	SetLevel(level)

	// 设置时间戳显示
	SetShowTime(cfg.ShowTimestamp)

	// 设置输出目标
	var output io.Writer
	switch cfg.Output {
	case "", "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		// 尝试打开文件
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file %s: %w", cfg.Output, err)
		}
		output = file
		// 注意：文件不会自动关闭，需要在程序退出时处理
	}
	SetOutput(output)

	return nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Level:         "info",
		Output:        "stdout",
		ShowTimestamp: false,
	}
}
