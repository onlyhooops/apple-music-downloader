package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogLevel 日志等级
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR"}

// String 返回日志等级的字符串表示
func (l LogLevel) String() string {
	if int(l) < len(levelNames) {
		return levelNames[l]
	}
	return "UNKNOWN"
}

// Logger 日志记录器接口
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	SetLevel(level LogLevel)
	SetOutput(w io.Writer)
	SetShowTime(show bool)
}

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	mu       sync.Mutex
	level    LogLevel
	output   io.Writer
	showTime bool
}

// New 创建一个新的日志记录器
func New() *DefaultLogger {
	return &DefaultLogger{
		level:    INFO,
		output:   os.Stdout,
		showTime: false, // UI模式下默认不显示时间戳
	}
}

// log 内部日志方法（带锁保护）
func (l *DefaultLogger) log(level LogLevel, format string, args ...interface{}) {
	// 日志等级过滤
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// 构建前缀
	var prefix string
	if l.showTime {
		prefix = fmt.Sprintf("[%s] %s: ",
			time.Now().Format("15:04:05"),
			levelNames[level])
	}

	// 格式化并输出
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "%s%s\n", prefix, message)
}

// Debug 输出DEBUG级别日志
func (l *DefaultLogger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 输出INFO级别日志
func (l *DefaultLogger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn 输出WARN级别日志
func (l *DefaultLogger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error 输出ERROR级别日志
func (l *DefaultLogger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// SetLevel 设置日志等级
func (l *DefaultLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput 设置输出目标
func (l *DefaultLogger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

// SetShowTime 设置是否显示时间戳
func (l *DefaultLogger) SetShowTime(show bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.showTime = show
}

// 全局logger实例
var global = New()

// Debug 全局Debug日志
func Debug(format string, args ...interface{}) {
	global.Debug(format, args...)
}

// Info 全局Info日志
func Info(format string, args ...interface{}) {
	global.Info(format, args...)
}

// Warn 全局Warn日志
func Warn(format string, args ...interface{}) {
	global.Warn(format, args...)
}

// Error 全局Error日志
func Error(format string, args ...interface{}) {
	global.Error(format, args...)
}

// SetLevel 设置全局日志等级
func SetLevel(level LogLevel) {
	global.SetLevel(level)
}

// SetOutput 设置全局输出目标
func SetOutput(w io.Writer) {
	global.SetOutput(w)
}

// SetShowTime 设置全局是否显示时间戳
func SetShowTime(show bool) {
	global.SetShowTime(show)
}

// ParseLevel 从字符串解析日志等级
func ParseLevel(s string) LogLevel {
	switch s {
	case "debug", "DEBUG":
		return DEBUG
	case "info", "INFO":
		return INFO
	case "warn", "WARN", "warning", "WARNING":
		return WARN
	case "error", "ERROR":
		return ERROR
	default:
		return INFO
	}
}
