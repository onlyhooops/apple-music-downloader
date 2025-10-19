package ui

import (
	"fmt"
	"main/internal/progress"

	"github.com/fatih/color"
)

// UIProgressListener UI进度监听器
// 实现progress.ProgressListener接口
// 将进度事件转换为UI更新
type UIProgressListener struct {
	// 可以添加需要的状态
}

// NewUIProgressListener 创建UI进度监听器
func NewUIProgressListener() *UIProgressListener {
	return &UIProgressListener{}
}

// OnProgress 处理进度更新事件
func (l *UIProgressListener) OnProgress(event progress.ProgressEvent) {
	status := formatStatus(event)
	colorFunc := getColorFunc(event.Stage)
	UpdateStatus(event.TrackIndex, status, colorFunc)
}

// OnComplete 处理完成事件
func (l *UIProgressListener) OnComplete(trackIndex int) {
	greenFunc := color.New(color.FgGreen).SprintFunc()
	UpdateStatus(trackIndex, "下载完成", greenFunc)
}

// OnError 处理错误事件
func (l *UIProgressListener) OnError(trackIndex int, err error) {
	errMsg := truncateError(err)
	redFunc := color.New(color.FgRed).SprintFunc()
	UpdateStatus(trackIndex, errMsg, redFunc)
}

// formatStatus 根据进度事件格式化状态文本
func formatStatus(event progress.ProgressEvent) string {
	// 如果事件已提供状态文本，直接使用
	if event.Status != "" {
		return event.Status
	}

	// 否则根据阶段和进度格式化
	switch event.Stage {
	case "download":
		if event.Percentage >= 0 {
			speedStr := formatSpeed(event.SpeedBPS)
			if speedStr != "" {
				return fmt.Sprintf("下载中 %d%% (%s)",
					event.Percentage, speedStr)
			}
			// 速度过低或为0，只显示百分比
			return fmt.Sprintf("下载中 %d%%", event.Percentage)
		}
		return "准备下载..."

	case "decrypt":
		if event.Percentage >= 0 {
			return fmt.Sprintf("解密中 %d%%", event.Percentage)
		}
		return "准备解密..."

	case "tag":
		return "写入标签中..."

	case "complete":
		return "下载完成"

	case "error":
		if event.Error != nil {
			return truncateError(event.Error)
		}
		return "发生错误"

	default:
		// 未知阶段，尝试从Status字段获取
		if event.Status != "" {
			return event.Status
		}
		return "处理中..."
	}
}

// getColorFunc 根据阶段返回对应的颜色函数
func getColorFunc(stage string) func(...interface{}) string {
	switch stage {
	case "download", "decrypt":
		return color.New(color.FgYellow).SprintFunc()
	case "tag":
		return color.New(color.FgCyan).SprintFunc()
	case "complete":
		return color.New(color.FgGreen).SprintFunc()
	case "error":
		return color.New(color.FgRed).SprintFunc()
	default:
		// 默认白色
		return func(a ...interface{}) string {
			return fmt.Sprint(a...)
		}
	}
}

// formatSpeed 格式化速度（字节/秒 → MB/s）
func formatSpeed(bps float64) string {
	// 如果速度过低（小于10KB/s），可能处于等待或缓冲状态，不显示速度
	if bps < 10240 {
		return ""
	}
	mbps := bps / 1024 / 1024
	// 对于小于0.1MB/s的速度，显示KB/s会更直观
	if mbps < 0.1 {
		kbps := bps / 1024
		return fmt.Sprintf("%.0f KB/s", kbps)
	}
	return fmt.Sprintf("%.1f MB/s", mbps)
}

// truncateError 截断错误信息到合适长度
// 避免错误信息过长导致UI显示混乱
func truncateError(err error) string {
	if err == nil {
		return ""
	}

	msg := err.Error()
	maxLength := 50

	// 根据终端宽度动态调整
	termWidth := getTerminalWidth()
	if termWidth > 60 {
		maxLength = 60
	} else if termWidth > 40 {
		maxLength = 40
	} else {
		maxLength = 30
	}

	if len(msg) > maxLength {
		return msg[:maxLength] + "..."
	}
	return msg
}
