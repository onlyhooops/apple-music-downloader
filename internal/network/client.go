package network

import (
	"net"
	"net/http"
	"time"

	"main/internal/logger"
	"main/utils/structs"
)

var (
	// DefaultClient 默认HTTP客户端（用于远程服务）
	DefaultClient *http.Client

	// LocalWrapperClient 本地wrapper服务专用HTTP客户端（启用优化）
	LocalWrapperClient *http.Client

	// 是否启用本地优化模式
	localOptimizationEnabled bool
)

// InitializeClients 初始化 HTTP 客户端
func InitializeClients(config *structs.ConfigSet) {
	localOptimizationEnabled = config.LocalWrapperOptimization.Enabled

	// 初始化默认客户端（用于远程 API 调用）
	DefaultClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// 根据配置初始化本地 wrapper 客户端
	if localOptimizationEnabled {
		initializeLocalWrapperClient(config.LocalWrapperOptimization)
	} else {
		// 未启用优化时，使用默认配置但仍创建独立客户端
		LocalWrapperClient = &http.Client{
			Timeout: 60 * time.Second, // wrapper 服务可能需要更长时间
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          50,
				MaxIdleConnsPerHost:   10,
				IdleConnTimeout:       90 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}
		logger.Info("🔌 本地 wrapper 优化: 未启用（使用默认配置）")
	}
}

// initializeLocalWrapperClient 初始化本地 wrapper 服务专用客户端（优化配置）
func initializeLocalWrapperClient(config structs.LocalWrapperConfig) {
	// 应用默认值
	maxIdleConns := config.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 200 // 本地服务可以支持更多连接
	}

	maxIdleConnsPerHost := config.MaxIdleConnsPerHost
	if maxIdleConnsPerHost == 0 {
		maxIdleConnsPerHost = 100 // 本地服务每个端口更多连接
	}

	maxConnsPerHost := config.MaxConnsPerHost
	if maxConnsPerHost == 0 {
		maxConnsPerHost = 0 // 0 表示不限制
	}

	idleConnTimeout := config.IdleConnTimeoutSec
	if idleConnTimeout == 0 {
		idleConnTimeout = 300 // 本地服务保持更长时间（5分钟）
	}

	dialTimeout := config.DialTimeoutMs
	if dialTimeout == 0 {
		dialTimeout = 100 // 本地连接应该很快（100ms）
	}

	keepAlive := config.KeepAlive
	if !keepAlive {
		keepAlive = true // 默认启用 KeepAlive
	}

	expectContinueTime := config.ExpectContinueTimeMs
	if expectContinueTime == 0 {
		expectContinueTime = 100 // 本地服务快速响应（100ms）
	}

	// 创建优化的 Transport
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(dialTimeout) * time.Millisecond,
			KeepAlive: 60 * time.Second, // 本地连接保持活跃
		}).DialContext,
		MaxIdleConns:          maxIdleConns,
		MaxIdleConnsPerHost:   maxIdleConnsPerHost,
		MaxConnsPerHost:       maxConnsPerHost,
		IdleConnTimeout:       time.Duration(idleConnTimeout) * time.Second,
		TLSHandshakeTimeout:   500 * time.Millisecond, // 本地通常不需要 TLS
		ExpectContinueTimeout: time.Duration(expectContinueTime) * time.Millisecond,
		DisableKeepAlives:     !keepAlive,
		DisableCompression:    config.DisableCompression, // 本地通讯可禁用压缩
		WriteBufferSize:       32 * 1024,                 // 32KB 写缓冲
		ReadBufferSize:        32 * 1024,                 // 32KB 读缓冲
	}

	LocalWrapperClient = &http.Client{
		Timeout:   120 * time.Second, // wrapper 服务可能需要更长时间
		Transport: transport,
	}

	logger.Info("🚀 本地 wrapper 优化: 已启用")
	logger.Debug("[网络优化] 最大空闲连接: %d", maxIdleConns)
	logger.Debug("[网络优化] 每主机最大空闲连接: %d", maxIdleConnsPerHost)
	logger.Debug("[网络优化] 每主机最大连接: %d (0=不限制)", maxConnsPerHost)
	logger.Debug("[网络优化] 空闲连接超时: %ds", idleConnTimeout)
	logger.Debug("[网络优化] 连接超时: %dms", dialTimeout)
	logger.Debug("[网络优化] TCP KeepAlive: %v", keepAlive)
	logger.Debug("[网络优化] 禁用压缩: %v", config.DisableCompression)
}

// GetWrapperClient 获取用于 wrapper 服务的 HTTP 客户端
func GetWrapperClient() *http.Client {
	if LocalWrapperClient != nil {
		return LocalWrapperClient
	}
	return DefaultClient
}

// IsLocalOptimizationEnabled 检查是否启用了本地优化
func IsLocalOptimizationEnabled() bool {
	return localOptimizationEnabled
}

// GetDefaultClient 获取默认 HTTP 客户端
func GetDefaultClient() *http.Client {
	if DefaultClient != nil {
		return DefaultClient
	}
	return http.DefaultClient
}
