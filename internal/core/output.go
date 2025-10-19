package core

import (
	"fmt"
	"main/internal/logger"
	"strings"
	"sync"
)

// OutputMutex 全局输出互斥锁，用于保护所有标准输出操作
// 防止多个goroutine的输出与动态UI渲染相互干扰
// 注意：UI渲染仍使用此锁，但日志输出已迁移到logger包的独立锁
var OutputMutex sync.Mutex

// SafePrintf 线程安全的Printf封装
// Deprecated: 使用 logger.Info() 替代。此函数保留用于向后兼容。
// 通过logger统一管理日志输出，支持日志等级控制和配置化输出
func SafePrintf(format string, a ...interface{}) {
	// 转发到logger系统
	// 移除格式字符串末尾的\n（logger会自动添加）
	format = strings.TrimSuffix(format, "\n")
	logger.Info(format, a...)
}

// SafePrintln 线程安全的Println封装
// Deprecated: 使用 logger.Info() 替代。此函数保留用于向后兼容。
// 通过logger统一管理日志输出，支持日志等级控制和配置化输出
func SafePrintln(a ...interface{}) {
	// 转发到logger系统
	msg := fmt.Sprint(a...)
	logger.Info("%s", msg)
}

// SafePrint 线程安全的Print封装
// Deprecated: 使用 logger.Info() 替代。此函数保留用于向后兼容。
// 通过logger统一管理日志输出，支持日志等级控制和配置化输出
func SafePrint(a ...interface{}) {
	// 转发到logger系统
	msg := fmt.Sprint(a...)
	logger.Info("%s", msg)
}
