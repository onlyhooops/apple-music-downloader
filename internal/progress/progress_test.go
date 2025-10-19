package progress

import (
	"sync"
	"testing"
	"time"
)

// MockListener 测试用的mock监听器
type MockListener struct {
	progressEvents []ProgressEvent
	completeEvents []int
	errorEvents    map[int]error
	mu             sync.Mutex
}

func NewMockListener() *MockListener {
	return &MockListener{
		progressEvents: make([]ProgressEvent, 0),
		completeEvents: make([]int, 0),
		errorEvents:    make(map[int]error),
	}
}

func (m *MockListener) OnProgress(event ProgressEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.progressEvents = append(m.progressEvents, event)
}

func (m *MockListener) OnComplete(trackIndex int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.completeEvents = append(m.completeEvents, trackIndex)
}

func (m *MockListener) OnError(trackIndex int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorEvents[trackIndex] = err
}

func (m *MockListener) GetProgressCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.progressEvents)
}

func (m *MockListener) GetCompleteCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.completeEvents)
}

// TestProgressNotifier 测试进度通知器基础功能
func TestProgressNotifier(t *testing.T) {
	notifier := NewNotifier()
	listener := NewMockListener()
	
	notifier.AddListener(listener)
	
	// 发送进度事件
	event := ProgressEvent{
		TrackIndex: 0,
		Stage:      "download",
		Percentage: 50,
		SpeedBPS:   1024000,
	}
	notifier.Notify(event)
	
	// 验证事件被接收
	time.Sleep(10 * time.Millisecond) // 给异步处理一点时间
	if listener.GetProgressCount() != 1 {
		t.Errorf("Expected 1 progress event, got %d", listener.GetProgressCount())
	}
}

// TestProgressNotifierMultipleListeners 测试多个监听器
func TestProgressNotifierMultipleListeners(t *testing.T) {
	notifier := NewNotifier()
	listener1 := NewMockListener()
	listener2 := NewMockListener()
	
	notifier.AddListener(listener1)
	notifier.AddListener(listener2)
	
	// 发送事件
	event := ProgressEvent{
		TrackIndex: 0,
		Stage:      "download",
		Percentage: 75,
	}
	notifier.Notify(event)
	
	// 验证两个监听器都收到
	time.Sleep(10 * time.Millisecond)
	if listener1.GetProgressCount() != 1 {
		t.Errorf("Listener1: expected 1 event, got %d", listener1.GetProgressCount())
	}
	if listener2.GetProgressCount() != 1 {
		t.Errorf("Listener2: expected 1 event, got %d", listener2.GetProgressCount())
	}
}

// TestProgressNotifierComplete 测试完成事件
func TestProgressNotifierComplete(t *testing.T) {
	notifier := NewNotifier()
	listener := NewMockListener()
	
	notifier.AddListener(listener)
	notifier.NotifyComplete(0)
	
	time.Sleep(10 * time.Millisecond)
	if listener.GetCompleteCount() != 1 {
		t.Errorf("Expected 1 complete event, got %d", listener.GetCompleteCount())
	}
}

// TestProgressNotifierConcurrency 测试并发安全
func TestProgressNotifierConcurrency(t *testing.T) {
	notifier := NewNotifier()
	listener := NewMockListener()
	
	notifier.AddListener(listener)
	
	var wg sync.WaitGroup
	numGoroutines := 100
	eventsPerGoroutine := 10
	
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				notifier.Notify(ProgressEvent{
					TrackIndex: id,
					Percentage: j * 10,
				})
			}
		}(i)
	}
	
	wg.Wait()
	time.Sleep(50 * time.Millisecond)
	
	expected := numGoroutines * eventsPerGoroutine
	actual := listener.GetProgressCount()
	if actual != expected {
		t.Errorf("Expected %d events, got %d", expected, actual)
	}
}

// TestProgressAdapter 测试适配器功能
func TestProgressAdapter(t *testing.T) {
	notifier := NewNotifier()
	listener := NewMockListener()
	notifier.AddListener(listener)
	
	adapter := NewProgressAdapter(notifier, 0, "download")
	progressChan := adapter.ToChan()
	
	// 发送旧格式的进度更新
	progressChan <- ProgressUpdate{
		Percentage: 50,
		SpeedBPS:   2048000,
		Stage:      "download",
	}
	
	progressChan <- ProgressUpdate{
		Percentage: 100,
		SpeedBPS:   2048000,
		Stage:      "download",
	}
	
	close(progressChan)
	
	time.Sleep(50 * time.Millisecond)
	
	// 验证事件被转换和接收
	if listener.GetProgressCount() != 2 {
		t.Errorf("Expected 2 events, got %d", listener.GetProgressCount())
	}
	
	// 验证事件内容
	listener.mu.Lock()
	events := listener.progressEvents
	listener.mu.Unlock()
	
	if len(events) >= 1 {
		if events[0].Percentage != 50 {
			t.Errorf("Expected percentage 50, got %d", events[0].Percentage)
		}
		if events[0].Stage != "download" {
			t.Errorf("Expected stage 'download', got %s", events[0].Stage)
		}
	}
}

// TestProgressAdapterStageUpdate 测试适配器阶段更新
func TestProgressAdapterStageUpdate(t *testing.T) {
	notifier := NewNotifier()
	listener := NewMockListener()
	notifier.AddListener(listener)
	
	adapter := NewProgressAdapter(notifier, 0, "download")
	progressChan := adapter.ToChan()
	
	// 发送下载进度
	progressChan <- ProgressUpdate{Percentage: 50, Stage: "download"}
	
	// 更新阶段
	adapter.UpdateStage("decrypt")
	
	// 发送解密进度
	progressChan <- ProgressUpdate{Percentage: 25, Stage: "decrypt"}
	
	close(progressChan)
	time.Sleep(50 * time.Millisecond)
	
	if listener.GetProgressCount() != 2 {
		t.Errorf("Expected 2 events, got %d", listener.GetProgressCount())
	}
}

// TestListenerCount 测试监听器计数
func TestListenerCount(t *testing.T) {
	notifier := NewNotifier()
	
	if notifier.ListenerCount() != 0 {
		t.Error("Initial listener count should be 0")
	}
	
	notifier.AddListener(NewMockListener())
	if notifier.ListenerCount() != 1 {
		t.Error("Listener count should be 1 after adding")
	}
	
	notifier.AddListener(NewMockListener())
	if notifier.ListenerCount() != 2 {
		t.Error("Listener count should be 2 after adding second")
	}
}

// TestRemoveListener 测试移除监听器
func TestRemoveListener(t *testing.T) {
	notifier := NewNotifier()
	listener := NewMockListener()
	
	notifier.AddListener(listener)
	if notifier.ListenerCount() != 1 {
		t.Error("Should have 1 listener")
	}
	
	notifier.RemoveListener(listener)
	if notifier.ListenerCount() != 0 {
		t.Error("Should have 0 listeners after removal")
	}
}

