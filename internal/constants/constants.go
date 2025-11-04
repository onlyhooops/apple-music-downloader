package constants

import "time"

// ==================== 重试配置 ====================
const (
	// MaxRetryAttempts 每个账号的最大重试次数
	MaxRetryAttempts = 3

	// RetryDelayMilliseconds 重试间隔（毫秒）
	RetryDelayMilliseconds = 1500

	// MaxConnectionRefusedRetries 连接被拒绝的最大重试次数
	MaxConnectionRefusedRetries = 3
)

// ==================== Token 验证 ====================
const (
	// MinTokenLength Token 的最小长度
	MinTokenLength = 50
)

// ==================== 网络配置 ====================
const (
	// DefaultHTTPTimeout HTTP 请求默认超时时间
	DefaultHTTPTimeout = 30 * time.Second

	// DefaultUserAgent 默认用户代理
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	// iTunesUserAgent iTunes 客户端 User-Agent
	iTunesUserAgent = "iTunes/12.11.3 (Windows; Microsoft Windows 10 x64 Professional Edition (Build 19041); x64) AppleWebKit/7611.1022.4001.1 (dt:2)"
)

// ==================== 缓存配置 ====================
const (
	// DefaultCachePath 默认缓存路径
	DefaultCachePath = "./Cache"

	// DefaultBatchSize 默认批次大小
	DefaultBatchSize = 20

	// CacheHashLength 缓存目录哈希长度
	CacheHashLength = 16
)

// ==================== 路径长度限制 ====================
const (
	// WindowsMaxPathLength Windows 系统的最大路径长度
	WindowsMaxPathLength = 255

	// UnixMaxPathLength Unix/Linux 系统的最大路径长度
	UnixMaxPathLength = 4096
)

// ==================== 下载配置 ====================
const (
	// DefaultCoverSize 默认封面尺寸
	DefaultCoverSize = "5000x5000"

	// DefaultChunkThreads 默认切片下载线程数
	DefaultChunkThreads = 30

	// DefaultLosslessThreads 默认无损格式下载线程数
	DefaultLosslessThreads = 5

	// DefaultAacThreads 默认 AAC 格式下载线程数
	DefaultAacThreads = 5

	// DefaultHiResThreads 默认 Hi-Res 格式下载线程数
	DefaultHiResThreads = 5

	// DefaultMVThreads 默认 MV 下载线程数
	DefaultMVThreads = 3
)

// ==================== 工作-休息循环 ====================
const (
	// DefaultWorkMinutes 默认工作时长（分钟）
	DefaultWorkMinutes = 30

	// DefaultRestMinutes 默认休息时长（分钟）
	DefaultRestMinutes = 2

	// RestTickerInterval 休息倒计时显示间隔
	RestTickerInterval = 30 * time.Second

	// CleanupWaitSeconds 清理等待时间（秒）
	CleanupWaitSeconds = 2
)

// ==================== 音质参数 ====================
var (
	// ValidAlacSampleRates 有效的 ALAC 采样率
	ValidAlacSampleRates = []int{44100, 48000, 96000, 192000}

	// ValidAtmosBitrates 有效的 Atmos 码率
	ValidAtmosBitrates = []int{2448, 2768}

	// ValidAacTypes 有效的 AAC 类型
	ValidAacTypes = []string{"aac-lc", "aac", "aac-binaural", "aac-downmix"}

	// ValidMVResolutions 有效的 MV 分辨率
	ValidMVResolutions = []int{480, 720, 1080, 2160}
)

// ==================== 文件相关 ====================
const (
	// TempFilePrefix 临时文件前缀
	TempFilePrefix = ".apple-music-tmp-"

	// DefaultFilePermission 默认文件权限
	DefaultFilePermission = 0644

	// DefaultDirPermission 默认目录权限
	DefaultDirPermission = 0755
)

// ==================== API 相关 ====================
const (
	// APIPageLimit API 分页查询每页数量
	APIPageLimit = 100

	// APIPageOffset API 分页偏移增量
	APIPageOffset = 100
)

// ==================== Buffer 配置 ====================
const (
	// DefaultNetworkReadBufferKB 默认网络读取缓冲区大小（KB）
	DefaultNetworkReadBufferKB = 4096

	// DefaultBufferSizeKB 默认 I/O 缓冲区大小（KB）
	DefaultBufferSizeKB = 4
)

// ==================== 显示相关 ====================
const (
	// VisualSeparatorLength 视觉分隔符长度
	VisualSeparatorLength = 60

	// BannerSeparatorLength 横幅分隔符长度
	BannerSeparatorLength = 80

	// ErrorMessageMaxLength 错误信息最大显示长度
	ErrorMessageMaxLength = 40
)

