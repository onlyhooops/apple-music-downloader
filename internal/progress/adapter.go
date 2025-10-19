package progress

import (
	"sync"

	"main/utils/runv14"
)

// ProgressUpdate 旧的进度更新结构（保持兼容runv14/runv3）
// 这个结构体用于向后兼容现有的channel-based进度更新机制
type ProgressUpdate struct {
	Percentage int     // 进度百分比 (0-100)
	SpeedBPS   float64 // 速度（字节/秒）
	Stage      string  // 阶段: download/decrypt
}

// ProgressAdapter 适配器
// 将旧的channel-based进度更新适配为新的事件驱动模式
// 这样可以在不修改现有downloader代码的情况下接入新系统
type ProgressAdapter struct {
	notifier   *ProgressNotifier
	trackIndex int
	stage      string
	mu         sync.RWMutex // 保护stage字段的并发访问
}

// NewProgressAdapter 创建一个新的进度适配器
// notifier: 进度通知器
// trackIndex: 曲目索引
// stage: 当前阶段（download/decrypt）
func NewProgressAdapter(notifier *ProgressNotifier, trackIndex int, stage string) *ProgressAdapter {
	return &ProgressAdapter{
		notifier:   notifier,
		trackIndex: trackIndex,
		stage:      stage,
	}
}

// ToChan 创建一个兼容旧代码的channel
// 返回一个ProgressUpdate channel，旧代码可以继续往这个channel发送进度
// 适配器会在后台goroutine中将这些更新转换为ProgressEvent并通知监听器
func (a *ProgressAdapter) ToChan() chan<- ProgressUpdate {
	ch := make(chan ProgressUpdate, 10)

	// 启动后台goroutine转换进度更新
	go func() {
		for update := range ch {
			// 读取当前stage（需要锁保护）
			a.mu.RLock()
			currentStage := a.stage
			a.mu.RUnlock()

			// 将旧格式转换为新格式
			event := ProgressEvent{
				TrackIndex: a.trackIndex,
				Stage:      currentStage,
				Percentage: update.Percentage,
				SpeedBPS:   update.SpeedBPS,
				Status:     "", // 由UI监听器格式化
				Error:      nil,
				Metadata:   nil,
			}

			// 通知所有监听器
			a.notifier.Notify(event)
		}
	}()

	return ch
}

// UpdateStage 更新适配器的当前阶段
// 用于当下载器切换阶段时（从download到decrypt）
// 线程安全，使用锁保护stage字段
func (a *ProgressAdapter) UpdateStage(newStage string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.stage = newStage
}

// ToRunv14Chan 创建一个兼容runv14.ProgressUpdate的channel
// 专门用于适配runv14包的ProgressUpdate类型
func (a *ProgressAdapter) ToRunv14Chan() chan<- runv14.ProgressUpdate {
	ch := make(chan runv14.ProgressUpdate, 10)

	// 启动后台goroutine转换进度更新
	go func() {
		for update := range ch {
			// 读取当前stage（需要锁保护）
			a.mu.RLock()
			currentStage := a.stage
			a.mu.RUnlock()

			// 将runv14格式转换为Progress事件格式
			event := ProgressEvent{
				TrackIndex: a.trackIndex,
				Stage:      currentStage,
				Percentage: update.Percentage,
				SpeedBPS:   update.SpeedBPS,
				Status:     "", // 由UI监听器格式化
				Error:      nil,
				Metadata:   nil,
			}

			// 通知所有监听器
			a.notifier.Notify(event)
		}
	}()

	return ch
}

// Close 关闭适配器（如果需要手动关闭channel）
// 注意：通常让channel在下载完成时自然关闭即可
func (a *ProgressAdapter) Close(ch chan<- ProgressUpdate) {
	close(ch)
}
