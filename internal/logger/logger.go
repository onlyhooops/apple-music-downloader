package logger

import (
	"fmt"
	"main/internal/core"
	"strings"
	"sync"

	"github.com/fatih/color"
)

// LogMutex 全局日志输出锁，保护所有终端输出
var LogMutex sync.Mutex

// LogBuffer UI模式下的日志缓冲区
var LogBuffer []string

// LogLevel 日志级别
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelSuccess
)

// 初始化
func init() {
	LogBuffer = make([]string, 0, 100)
}

// 检查是否处于UI模式
func isUIMode() bool {
	core.UiMutex.Lock()
	defer core.UiMutex.Unlock()
	return len(core.TrackStatuses) > 0
}

// print 内部输出函数，处理UI模式检测和锁控制
func print(message string) {
	LogMutex.Lock()
	defer LogMutex.Unlock()

	// 如果处于UI模式，缓存日志
	if isUIMode() {
		LogBuffer = append(LogBuffer, message)
		return
	}

	// 否则直接输出
	fmt.Print(message)
}

// FlushBuffer 刷新缓冲区（UI结束后调用）
func FlushBuffer() {
	LogMutex.Lock()
	defer LogMutex.Unlock()

	if len(LogBuffer) > 0 {
		for _, msg := range LogBuffer {
			fmt.Print(msg)
		}
		LogBuffer = LogBuffer[:0] // 清空缓冲区
	}
}

// ClearBuffer 清空缓冲区（不输出）
func ClearBuffer() {
	LogMutex.Lock()
	defer LogMutex.Unlock()
	LogBuffer = LogBuffer[:0]
}

// Section 输出分隔区块标题（按照UI微调.txt规范）
func Section(title string) {
	separator := strings.Repeat("-", 50)
	print(separator + "\n")
	print(fmt.Sprintf("⏳ %s\n\n", title))
}

// SectionEnd 输出区块结束分隔线
func SectionEnd() {
	separator := strings.Repeat("-", 50)
	print(separator + "\n")
}

// Info 输出普通信息（绿色）
func Info(format string, args ...interface{}) {
	green := color.New(color.FgGreen).SprintFunc()
	message := fmt.Sprintf(format, args...)
	print(green(message) + "\n")
}

// Warn 输出警告信息（黄色，带警告emoji）
func Warn(format string, args ...interface{}) {
	yellow := color.New(color.FgYellow).SprintFunc()
	message := fmt.Sprintf(format, args...)
	print(fmt.Sprintf("⚠️  %s\n", yellow(message)))
}

// Error 输出错误信息（红色，带错误标记）
func Error(format string, args ...interface{}) {
	red := color.New(color.FgRed).SprintFunc()
	message := fmt.Sprintf(format, args...)
	print(fmt.Sprintf("❌ %s\n", red(message)))
}

// Success 输出成功信息（绿色，带成功emoji）
func Success(format string, args ...interface{}) {
	green := color.New(color.FgGreen).SprintFunc()
	message := fmt.Sprintf(format, args...)
	print(fmt.Sprintf("✅ %s\n", green(message)))
}

// Debug 输出调试信息（青色，仅在调试模式）
func Debug(format string, args ...interface{}) {
	if !core.Debug_mode {
		return
	}
	cyan := color.New(color.FgCyan).SprintFunc()
	message := fmt.Sprintf(format, args...)
	print(fmt.Sprintf("🔍 [DEBUG] %s\n", cyan(message)))
}

// Plain 输出无格式化的普通文本
func Plain(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	print(message)
}

// Printf 兼容fmt.Printf的接口（用于渐进式迁移）
func Printf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	print(message)
}

// Println 兼容fmt.Println的接口
func Println(args ...interface{}) {
	message := fmt.Sprintln(args...)
	print(message)
}

// AlbumHeader 专辑信息头部（按UI微调规范）
func AlbumHeader(artist, album string) {
	print(fmt.Sprintf("🎤 歌手: %s\n", artist))
	print(fmt.Sprintf("💽 专辑: %s\n", album))
}

// QualityInfo 音质信息（按UI微调规范）
func QualityInfo(quality string, threads int, regions string, accountCount int) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	
	print(fmt.Sprintf("🎵 %s: %s | %s | %s | %s\n",
		green("音源"),
		green(quality),
		green(fmt.Sprintf("%d 线程", threads)),
		yellow(regions),
		green(fmt.Sprintf("%d 个账户并行下载", accountCount)),
	))
}

// TaskProgress 任务进度（批量下载）
func TaskProgress(current, total int, message string) {
	if total > 1 {
		print(fmt.Sprintf("[%d/%d] %s\n", current, total, message))
	}
}

// Summary 输出任务总结（按UI微调规范）
func Summary(albumName, quality string, success, warn, errorCount int) {
	Success(fmt.Sprintf("下载完成"))
	print(fmt.Sprintf("💽 %s - %s\n\n", albumName, quality))
	print(fmt.Sprintf("已完成: %d  |  警告: %d  |  错误: %d\n", success, warn, errorCount))
	SectionEnd()
}

// TotalCount 输出总数和并发数
func TotalCount(total, concurrency int) {
	print(fmt.Sprintf("🔖 总数: %d, 并发数: %d\n", total, concurrency))
}

// SafePrint 线程安全的直接输出（用于特殊情况）
func SafePrint(text string) {
	LogMutex.Lock()
	defer LogMutex.Unlock()
	fmt.Print(text)
}

