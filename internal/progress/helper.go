package progress

// Helper 进度辅助函数
// 提供简化的API用于快速发送进度事件

// NotifyDownloadProgress 发送下载进度事件
func (n *ProgressNotifier) NotifyDownloadProgress(trackIndex int, percentage int, speedBPS float64) {
	n.Notify(ProgressEvent{
		TrackIndex: trackIndex,
		Stage:      "download",
		Percentage: percentage,
		SpeedBPS:   speedBPS,
	})
}

// NotifyDecryptProgress 发送解密进度事件
func (n *ProgressNotifier) NotifyDecryptProgress(trackIndex int, percentage int) {
	n.Notify(ProgressEvent{
		TrackIndex: trackIndex,
		Stage:      "decrypt",
		Percentage: percentage,
	})
}

// NotifyTag 发送标签写入事件
func (n *ProgressNotifier) NotifyTag(trackIndex int) {
	n.Notify(ProgressEvent{
		TrackIndex: trackIndex,
		Stage:      "tag",
		Status:     "写入标签中...",
	})
}

// NotifyStatus 发送自定义状态事件
func (n *ProgressNotifier) NotifyStatus(trackIndex int, status string, stage string) {
	n.Notify(ProgressEvent{
		TrackIndex: trackIndex,
		Stage:      stage,
		Status:     status,
	})
}
