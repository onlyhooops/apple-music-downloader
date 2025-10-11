package logger

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"
)

// TestLoggerLevel 测试日志等级过滤
func TestLoggerLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New()
	logger.SetOutput(buf)
	logger.SetLevel(WARN)

	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Warn("warn msg")
	logger.Error("error msg")

	output := buf.String()

	if strings.Contains(output, "debug msg") {
		t.Error("DEBUG message should be filtered")
	}
	if strings.Contains(output, "info msg") {
		t.Error("INFO message should be filtered")
	}
	if !strings.Contains(output, "warn msg") {
		t.Error("WARN message should be logged")
	}
	if !strings.Contains(output, "error msg") {
		t.Error("ERROR message should be logged")
	}
}

// TestLoggerOutput 测试输出重定向
func TestLoggerOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New()
	logger.SetOutput(buf)
	logger.SetLevel(INFO)

	testMsg := "test message"
	logger.Info("%s", testMsg)

	output := buf.String()
	if !strings.Contains(output, testMsg) {
		t.Errorf("Expected output to contain %q, got %q", testMsg, output)
	}
}

// TestLoggerFormat 测试格式化输出
func TestLoggerFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New()
	logger.SetOutput(buf)
	logger.SetLevel(INFO)

	format := "test %s %d"
	logger.Info(format, "message", 123)

	output := buf.String()
	if !strings.Contains(output, "test message 123") {
		t.Errorf("Format failed, got: %q", output)
	}
}

// TestLoggerConcurrency 测试并发安全
func TestLoggerConcurrency(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New()
	logger.SetOutput(buf)
	logger.SetLevel(INFO)

	var wg sync.WaitGroup
	numGoroutines := 100
	messagesPerGoroutine := 10

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Info("goroutine %d message %d", id, j)
			}
		}(i)
	}

	wg.Wait()

	// 验证输出行数
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	expected := numGoroutines * messagesPerGoroutine
	if len(lines) != expected {
		t.Errorf("Expected %d lines, got %d", expected, len(lines))
	}
}

// TestLoggerShowTime 测试时间戳显示
func TestLoggerShowTime(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New()
	logger.SetOutput(buf)
	logger.SetLevel(INFO)

	// 不显示时间戳
	logger.SetShowTime(false)
	logger.Info("without timestamp")
	output1 := buf.String()
	if strings.Contains(output1, "[") && strings.Contains(output1, "]") {
		t.Error("Should not contain timestamp")
	}

	// 显示时间戳
	buf.Reset()
	logger.SetShowTime(true)
	logger.Info("with timestamp")
	output2 := buf.String()
	if !strings.Contains(output2, "INFO:") {
		t.Error("Should contain log level with timestamp")
	}
}

// TestParseLevel 测试日志等级解析
func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", DEBUG},
		{"DEBUG", DEBUG},
		{"info", INFO},
		{"INFO", INFO},
		{"warn", WARN},
		{"WARN", WARN},
		{"warning", WARN},
		{"error", ERROR},
		{"ERROR", ERROR},
		{"unknown", INFO}, // 默认返回INFO
	}

	for _, tt := range tests {
		result := ParseLevel(tt.input)
		if result != tt.expected {
			t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

// TestGlobalLogger 测试全局logger
func TestGlobalLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	SetLevel(INFO)

	Info("global test")

	output := buf.String()
	if !strings.Contains(output, "global test") {
		t.Error("Global logger not working")
	}

	// 恢复默认输出
	SetOutput(io.Discard)
}

// TestLogLevelString 测试日志等级字符串表示
func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
	}

	for _, tt := range tests {
		result := tt.level.String()
		if result != tt.expected {
			t.Errorf("LogLevel(%d).String() = %q, want %q", tt.level, result, tt.expected)
		}
	}
}

