package logger

import (
	"io"
	"testing"
)

// BenchmarkLoggerInfo 测试Info日志性能
func BenchmarkLoggerInfo(b *testing.B) {
	logger := New()
	logger.SetOutput(io.Discard) // 输出到/dev/null，避免IO影响
	logger.SetLevel(INFO)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message %d", i)
	}
}

// BenchmarkLoggerInfoFiltered 测试被过滤的日志性能
func BenchmarkLoggerInfoFiltered(b *testing.B) {
	logger := New()
	logger.SetOutput(io.Discard)
	logger.SetLevel(ERROR) // 设置为ERROR，INFO会被过滤

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message %d", i)
	}
}

// BenchmarkLoggerError 测试Error日志性能
func BenchmarkLoggerError(b *testing.B) {
	logger := New()
	logger.SetOutput(io.Discard)
	logger.SetLevel(INFO)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("test error %d", i)
	}
}

// BenchmarkLoggerConcurrent 测试并发日志性能
func BenchmarkLoggerConcurrent(b *testing.B) {
	logger := New()
	logger.SetOutput(io.Discard)
	logger.SetLevel(INFO)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("concurrent test")
		}
	})
}

// BenchmarkLoggerWithTimestamp 测试带时间戳的日志性能
func BenchmarkLoggerWithTimestamp(b *testing.B) {
	logger := New()
	logger.SetOutput(io.Discard)
	logger.SetLevel(INFO)
	logger.SetShowTime(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test message %d", i)
	}
}

// BenchmarkGlobalLogger 测试全局logger性能
func BenchmarkGlobalLogger(b *testing.B) {
	SetOutput(io.Discard)
	SetLevel(INFO)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("test message %d", i)
	}
}

// BenchmarkParseLevel 测试日志等级解析性能
func BenchmarkParseLevel(b *testing.B) {
	levels := []string{"debug", "info", "warn", "error", "DEBUG", "INFO", "WARN", "ERROR"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseLevel(levels[i%len(levels)])
	}
}

