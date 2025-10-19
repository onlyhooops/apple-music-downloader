package progress

import (
	"sync"
)

// ProgressEvent 进度事件
type ProgressEvent struct {
	TrackIndex int                    // 曲目索引（在批次中）
	Stage      string                 // 阶段: download/decrypt/tag/complete/error
	Percentage int                    // 进度百分比 (0-100)
	SpeedBPS   float64                // 速度（字节/秒）
	Status     string                 // 状态描述文本
	Error      error                  // 错误信息（如有）
	Metadata   map[string]interface{} // 额外元数据
}

// ProgressListener 进度监听器接口
type ProgressListener interface {
	OnProgress(event ProgressEvent)
	OnComplete(trackIndex int)
	OnError(trackIndex int, err error)
}

// ProgressNotifier 进度通知器
// 实现观察者模式，管理多个监听器并分发进度事件
type ProgressNotifier struct {
	listeners []ProgressListener
	mu        sync.RWMutex
}

// NewNotifier 创建一个新的进度通知器
func NewNotifier() *ProgressNotifier {
	return &ProgressNotifier{
		listeners: make([]ProgressListener, 0),
	}
}

// AddListener 添加一个监听器
// 线程安全，可在运行时动态添加监听器
func (n *ProgressNotifier) AddListener(l ProgressListener) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.listeners = append(n.listeners, l)
}

// RemoveListener 移除一个监听器（可选功能）
func (n *ProgressNotifier) RemoveListener(l ProgressListener) {
	n.mu.Lock()
	defer n.mu.Unlock()
	
	for i, listener := range n.listeners {
		if listener == l {
			n.listeners = append(n.listeners[:i], n.listeners[i+1:]...)
			return
		}
	}
}

// Notify 通知所有监听器进度更新
// 使用RWMutex允许并发通知
func (n *ProgressNotifier) Notify(event ProgressEvent) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	
	for _, listener := range n.listeners {
		listener.OnProgress(event)
	}
}

// NotifyComplete 通知下载完成
func (n *ProgressNotifier) NotifyComplete(trackIndex int) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	
	for _, listener := range n.listeners {
		listener.OnComplete(trackIndex)
	}
}

// NotifyError 通知错误
func (n *ProgressNotifier) NotifyError(trackIndex int, err error) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	
	for _, listener := range n.listeners {
		listener.OnError(trackIndex, err)
	}
}

// ListenerCount 返回当前监听器数量（用于调试）
func (n *ProgressNotifier) ListenerCount() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return len(n.listeners)
}

